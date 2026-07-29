const getUserCode = () => {
  return "xxxxxxyxxxxxxxx".replace(/[xy]/g, function (c) {
    const r = (Math.random() * 16) | 0;
    const v = c === "x" ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
};

const block2DelayTime = 5000; // 5s
const block1DelayTime = 3000; // 3s

const focusDelayTime = 500; // 500ms
const blankDelayTime = 1000; // 500ms

const block2ShowTime = 300;
const block2MaskTime = 700;

let timer = null; //超时定时器回调

// 初始化jsPsych
const jsPsych = initJsPsych();

// let group_id = jsPsych.data.getURLVariable("group_id");
// let sample_id = jsPsych.data.getURLVariable("sample_id");

jsPsych.data.addProperties({
  userCode: getUserCode(),
  // group_id: group_id,
  // sample_id: sample_id,
});

const rand = jsPsych.randomization.randomInt(0, 1);

const resourceMap = getResource(jsPsych);

// console.log("resourceMap",resourceMap)

// 创建实验时间线
const timeline = [];

const preload = {
  type: jsPsychPreload,
  images: resourceMap.block1.list.map((c) => c.link),
  message: "资源加载中，请稍候...",
  error_message: "资源加载失败，请刷新页面重试",
};
timeline.push(preload);

// 欢迎界面
const welcome_trial = {
  type: jsPsychHtmlKeyboardResponse,
  stimulus: `
          <div style="text-align: center;">
              <h2>素材测试</h2>
              <p>玩家您好！</p>
              <p>本小程序会呈现一些美术素材，</p>
              <p>请根据您的第一印象，尽量在3秒内进行选择。</p>
              <p>您的选择没有对错之分，请按照您的真实想法选择即可</p>
              <button class="btn btn-yes" id="start">确认</button>
          </div>
      `,
  choices: "NO_KEYS",
  on_load: function () {
    document.getElementById("start").addEventListener("click", function () {
      jsPsych.finishTrial();
    });
  },
};
timeline.push(welcome_trial);

timeline.push({
  type: jsPsychPreload,
  images: resourceMap.block2.list.map((c) => c.link),
  message: "资源加载中，请稍候...",
  error_message: "资源加载失败，请刷新页面重试",
});

// 回合1
injectBlock1(timeline, rand, resourceMap.block1);

let preload2 = {
  type: jsPsychPreload,
  images: resourceMap.block2.list.map((c) => c.link),
  message: "资源加载中，请稍候...",
  error_message: "资源加载失败，请刷新页面重试",
};

timeline.push({
  timeline: [preload2],
  conditional_function: function () {
    const ch = jsPsych.data
      .get()
      .filter({ block: 1, 图片: "测谎图" })
      .last(1)
      .values()[0];
    if (ch && ch.response === "no") {
      return false;
    }

    return true;
  },
});

// 回合2
injectBlock2(timeline, rand, resourceMap.block2, resource.mask);

// 开始实验
jsPsych.run(timeline);
