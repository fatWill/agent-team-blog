(function (global) {
  'use strict';

  var SDK_NAME = 'wj-jspsych-sdk';
  var SDK_VERSION = '7.0.0';
  /** 本 SDK 适配的 jsPsych 主版本 */
  var TARGET_MAJOR = 7;
  /** jsPsych 7.x 用于提前结束实验的方法 */
  var ABORT_METHODS = ['endExperiment', 'endCurrentTimeline'];

  /** 与问卷 HTML 题一致的通信关键词，勿改 */
  var END_KEY = '__WJ_IFRAME_QUESTION_END__';
  var INIT_KEY = '__WJ_IFRAME_INIT_DATA__';

  /** 单个字段最大字符数，超出即丢弃，避免把 base64 素材塞进 postMessage */
  var MAX_VALUE_LENGTH = 20000;
  /** 单条消息最大字符数 */
  var MAX_MESSAGE_LENGTH = 200000;

  var LOG_PREFIX = '[' + SDK_NAME + '@' + SDK_VERSION + '] ';

  /* ------------------------------------------------------------------ *
   * 基础工具
   * ------------------------------------------------------------------ */

  function hasConsole(method) {
    return !!(global.console && global.console[method]);
  }

  function logInfo(message) {
    if (hasConsole('log')) {
      global.console.log(LOG_PREFIX + message);
    }
  }

  function logWarn(message) {
    if (hasConsole('warn')) {
      global.console.warn(LOG_PREFIX + message);
    }
  }

  function logError(message) {
    if (hasConsole('error')) {
      global.console.error(LOG_PREFIX + message);
    }
  }

  function isFn(value) {
    return typeof value === 'function';
  }

  function isPlainObject(value) {
    return !!value && typeof value === 'object' && !isArray(value);
  }

  function isArray(value) {
    return Object.prototype.toString.call(value) === '[object Array]';
  }

  /** 安全调用用户回调：用户代码异常不应中断上报 */
  function safeCall(fn, thisArg, args) {
    if (!isFn(fn)) {
      return undefined;
    }
    try {
      return fn.apply(thisArg, args || []);
    } catch (e) {
      logError('实验回调执行异常：' + (e && e.message ? e.message : e));
      return undefined;
    }
  }

  /* ------------------------------------------------------------------ *
   * 环境校验
   * ------------------------------------------------------------------ */

  /** 是否运行在 iframe 中 */
  function inIframe() {
    try {
      return global.self !== global.top;
    } catch (e) {
      // 跨域访问 window.top 抛错，说明确实被嵌套在别的文档里
      return true;
    }
  }

  /** 在页面上渲染一个不可关闭的错误提示 */
  function renderFatal(message) {
    var render = function () {
      var body = global.document && global.document.body;
      if (!body) {
        return;
      }
      var mask = global.document.createElement('div');
      mask.setAttribute('data-wj-sdk-fatal', '1');
      mask.style.cssText = [
        'position:fixed',
        'top:0',
        'left:0',
        'right:0',
        'bottom:0',
        'z-index:2147483647',
        'display:flex',
        'align-items:center',
        'justify-content:center',
        'padding:24px',
        'box-sizing:border-box',
        'background:#fff',
        'color:#d54941',
        'font-size:16px',
        'line-height:1.6',
        'text-align:center',
        'font-family:-apple-system,BlinkMacSystemFont,"Helvetica Neue",Arial,sans-serif'
      ].join(';');
      mask.appendChild(global.document.createTextNode(message));
      body.appendChild(mask);
    };

    if (global.document && global.document.body) {
      render();
    } else if (global.document) {
      global.document.addEventListener('DOMContentLoaded', render);
    }
  }

  /* ------------------------------------------------------------------ *
   * 与父页面（问卷答题页）通信
   * ------------------------------------------------------------------ */

  /** 父页面在 iframe onload 后下发的历史作答数据 */
  var initData = [];
  /** 是否已上报过结束信号 */
  var finished = false;

  /**
   * 推断父页面 origin：
   * 优先用 ancestorOrigins（Chromium/WebKit 支持），其次用 referrer，
   * 都拿不到时退化为 '*'（父页面本身会校验来源，功能不受影响）。
   * 也支持实验页面通过 window.__WJ_PARENT_ORIGIN__ 显式指定。
   */
  function getParentOrigin() {
    if (typeof global.__WJ_PARENT_ORIGIN__ === 'string' && global.__WJ_PARENT_ORIGIN__) {
      return global.__WJ_PARENT_ORIGIN__;
    }
    try {
      var ancestors = global.location.ancestorOrigins;
      if (ancestors && ancestors.length) {
        return ancestors[ancestors.length - 1];
      }
    } catch (e) {
      // 忽略，继续走 referrer
    }
    try {
      var referrer = global.document.referrer;
      if (referrer) {
        var matched = /^(https?:\/\/[^/?#]+)/i.exec(referrer);
        if (matched) {
          return matched[1];
        }
      }
    } catch (e) {
      // 忽略
    }
    return '*';
  }

  /** 把值裁剪成可安全传输的形态：过长或无法序列化的一律丢弃 */
  function normalizeValue(value) {
    if (value === null || value === undefined) {
      return value;
    }
    var type = typeof value;
    if (type === 'number' || type === 'boolean' || type === 'string') {
      if (type === 'string' && value.length > MAX_VALUE_LENGTH) {
        return undefined;
      }
      return value;
    }
    if (type === 'function') {
      return undefined;
    }
    try {
      var serialized = JSON.stringify(value);
      if (serialized === undefined || serialized.length > MAX_VALUE_LENGTH) {
        return undefined;
      }
      return JSON.parse(serialized);
    } catch (e) {
      return undefined;
    }
  }

  /** 把一条 trial 数据转成可上报的扁平对象（键即问卷的自定义变量名） */
  function normalizeTrialData(data) {
    var payload = {};
    if (!isPlainObject(data)) {
      return payload;
    }
    var keys = Object.keys(data);
    for (var i = 0; i < keys.length; i++) {
      var key = keys[i];
      var value = normalizeValue(data[key]);
      if (value !== undefined) {
        payload[key] = value;
      }
    }
    return payload;
  }

  /** 发送一条消息给父页面 */
  function post(payload) {
    if (!isPlainObject(payload) || !Object.keys(payload).length) {
      return;
    }
    var serialized;
    try {
      serialized = JSON.stringify(payload);
    } catch (e) {
      logWarn('数据无法序列化，已跳过本次上报');
      return;
    }
    if (serialized.length > MAX_MESSAGE_LENGTH) {
      logWarn('单条数据超过 ' + MAX_MESSAGE_LENGTH + ' 字符，已跳过本次上报');
      return;
    }
    try {
      global.parent.postMessage(payload, getParentOrigin());
    } catch (e) {
      logError('上报失败：' + (e && e.message ? e.message : e));
    }
  }

  /** 上报一条 trial 数据 */
  function reportTrial(data) {
    post(normalizeTrialData(data));
  }

  /** 上报结束信号（幂等，重复调用只生效一次） */
  function reportEnd() {
    if (finished) {
      return;
    }
    finished = true;
    var payload = {};
    payload[END_KEY] = true;
    post(payload);
    logInfo('实验已结束，已通知问卷本题作答完毕');
  }

  /** 监听父页面下发的历史作答数据 */
  function listenParent() {
    if (!global.addEventListener) {
      return;
    }
    global.addEventListener('message', function (event) {
      var data = event && event.data;
      if (typeof data === 'string') {
        try {
          data = JSON.parse(data);
        } catch (e) {
          return;
        }
      }
      if (!isPlainObject(data) || !(INIT_KEY in data)) {
        return;
      }
      initData = isArray(data[INIT_KEY]) ? data[INIT_KEY] : [];
      logInfo('已接收问卷下发的历史作答数据，共 ' + initData.length + ' 批');
    });
  }

  /* ------------------------------------------------------------------ *
   * 上报去重
   *
   * jsPsych 每结束一个 trial 会先触发 on_data_update（写入数据），
   * 再触发 trial 自身的 on_finish、最后触发全局 on_trial_finish。
   * trial 的 on_finish 可能会改写数据，所以以 on_trial_finish 为准；
   * on_data_update 延迟一拍再判断，只用于兜底那些不走 trial 流程的数据写入。
   * ------------------------------------------------------------------ */

  /** 已上报过的 trial_index 集合 */
  var reportedIndexes = {};

  function markReported(data) {
    if (isPlainObject(data) && data.trial_index !== undefined && data.trial_index !== null) {
      reportedIndexes[String(data.trial_index)] = true;
    }
  }

  function isReported(data) {
    if (!isPlainObject(data) || data.trial_index === undefined || data.trial_index === null) {
      return false;
    }
    return reportedIndexes[String(data.trial_index)] === true;
  }

  function reportTrialOnce(data) {
    if (isReported(data)) {
      return;
    }
    markReported(data);
    reportTrial(data);
  }

  /* ------------------------------------------------------------------ *
   * 接管 initJsPsych 的配置与实例
   * ------------------------------------------------------------------ */

  function checkVersion(instance) {
    if (!instance || !isFn(instance.version)) {
      return;
    }
    var version = safeCall(instance.version, instance);
    if (typeof version !== 'string') {
      return;
    }
    var major = parseInt(version.split('.')[0], 10);
    if (major !== TARGET_MAJOR) {
      logWarn(
        '当前 jsPsych 版本为 ' +
          version +
          '，本 SDK 适配 jsPsych ' +
          TARGET_MAJOR +
          '.x，请改用 ' +
          SDK_NAME +
          '@' +
          major +
          '.0.0.js'
      );
    }
  }

  /** 包装 initJsPsych 的配置项，注入数据上报 */
  function wrapOptions(options) {
    var opts = isPlainObject(options) ? options : {};

    var userOnTrialFinish = opts.on_trial_finish;
    var userOnDataUpdate = opts.on_data_update;
    var userOnFinish = opts.on_finish;
    var userOnClose = opts.on_close;

    opts.on_trial_finish = function (data) {
      reportTrialOnce(data);
      return safeCall(userOnTrialFinish, this, [data]);
    };

    opts.on_data_update = function (data) {
      // 延迟一拍：若 on_trial_finish 已上报同一 trial，则不重复上报
      var snapshot = data;
      global.setTimeout(function () {
        reportTrialOnce(snapshot);
      }, 0);
      return safeCall(userOnDataUpdate, this, [data]);
    };

    opts.on_finish = function (data) {
      var result = safeCall(userOnFinish, this, [data]);
      reportEnd();
      return result;
    };

    opts.on_close = function () {
      // 中途关闭/刷新不代表作答完成，只保证已产生的数据都已上报，不发结束信号
      return safeCall(userOnClose, this, arguments);
    };

    return opts;
  }

  /** 包装 jsPsych 实例：run 结束与提前中断都要保证发出结束信号 */
  function wrapInstance(instance) {
    if (!instance || instance.__wjAttached) {
      return instance;
    }
    instance.__wjAttached = true;
    checkVersion(instance);

    for (var i = 0; i < ABORT_METHODS.length; i++) {
      (function (name) {
        var origin = instance[name];
        if (!isFn(origin)) {
          return;
        }
        instance[name] = function () {
          var result = origin.apply(this, arguments);
          // on_finish 正常会随后触发；这里兜底，reportEnd 幂等不会重复发送
          global.setTimeout(reportEnd, 0);
          return result;
        };
      })(ABORT_METHODS[i]);
    }

    // jsPsych 7 的 run 返回 Promise，resolve 即代表时间线跑完
    if (isFn(instance.run)) {
      var originRun = instance.run;
      instance.run = function () {
        var result = originRun.apply(this, arguments);
        if (result && isFn(result.then)) {
          result.then(reportEnd, function () {
            reportEnd();
          });
        }
        return result;
      };
    }

    return instance;
  }

  /** 生成被接管后的 initJsPsych */
  function wrapInitJsPsych(originInitJsPsych) {
    if (!isFn(originInitJsPsych) || originInitJsPsych.__wjWrapped) {
      return originInitJsPsych;
    }
    var wrapped = function (options) {
      if (!inIframe()) {
        var message = '本实验需在腾讯问卷内作答，请通过问卷链接进入。';
        renderFatal(message);
        throw new Error(LOG_PREFIX + message);
      }
      var instance = originInitJsPsych.call(this, wrapOptions(options));
      logInfo('已接管 jsPsych 数据上报');
      return wrapInstance(instance);
    };
    wrapped.__wjWrapped = true;
    return wrapped;
  }

  /**
   * 用访问器劫持全局属性，这样 SDK 与 jsPsych 的加载顺序都能兼容：
   * SDK 先加载则拦截后续赋值，jsPsych 先加载则立即包装。
   */
  function hookGlobalProperty(name, transform) {
    var value = transform(global[name]);
    try {
      Object.defineProperty(global, name, {
        configurable: true,
        enumerable: true,
        get: function () {
          return value;
        },
        set: function (next) {
          value = transform(next);
        }
      });
      return true;
    } catch (e) {
      // 属性不可重定义（极少见）：只能包装已加载的那份
      try {
        global[name] = value;
      } catch (e2) {
        logWarn('无法接管 ' + name + '，请改为手动调用 WJJsPsych.attach(jsPsych)');
        return false;
      }
      return true;
    }
  }

  function install() {
    if (global.WJJsPsych && global.WJJsPsych.__installed) {
      logWarn('SDK 已加载，忽略重复引入');
      return;
    }

    if (!inIframe()) {
      var message = '本实验需在腾讯问卷内作答，请通过问卷链接进入。';
      logError(message);
      renderFatal(message);
    }

    listenParent();

    // 全局函数形式：window.initJsPsych
    hookGlobalProperty('initJsPsych', function (value) {
      return isFn(value) ? wrapInitJsPsych(value) : value;
    });

    // 命名空间形式：window.jsPsychModule.initJsPsych
    hookGlobalProperty('jsPsychModule', function (value) {
      if (isPlainObject(value) && isFn(value.initJsPsych)) {
        value.initJsPsych = wrapInitJsPsych(value.initJsPsych);
      }
      return value;
    });

    global.WJJsPsych = {
      __installed: true,
      name: SDK_NAME,
      version: SDK_VERSION,
      targetMajor: TARGET_MAJOR,

      /**
       * 手动接管一个 jsPsych 实例。
       * 仅在实验用 ESM 打包（initJsPsych 不挂在 window 上）时需要：
       *   const jsPsych = initJsPsych(WJJsPsych.options({ ... }));
       *   WJJsPsych.attach(jsPsych);
       */
      attach: function (instance) {
        return wrapInstance(instance);
      },

      /** 供 ESM 场景手动包装 initJsPsych 的配置 */
      options: function (options) {
        return wrapOptions(options);
      },

      /** 父页面下发的历史作答数据（二维批次数组），用于续答场景 */
      getInitData: function () {
        return initData;
      },

      /** 手动上报一条数据（键需与问卷的自定义变量同名） */
      send: function (data) {
        reportTrial(data);
      },

      /** 手动上报结束信号（幂等） */
      finish: function () {
        reportEnd();
      }
    };

    logInfo('SDK 已就绪');
  }

  install();
})(typeof window !== 'undefined' ? window : this);
