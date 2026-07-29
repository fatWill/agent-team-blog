/**
 * 注入回合2
 * @param {*} timeline
 * @param {*} rand
 * @param {
 *  list: any[],
 *  last: any[],
 * } item
 */
function injectBlock2(timeline, rand, item, mask) {
  let temp = []
  const welcome_trial = {
    type: jsPsychHtmlKeyboardResponse,
    stimulus: `
    <p>接下来在屏幕中会出现一张美术素材，随后其会迅速消失；请判断该素材是否在刚才出现过</p>
    <p>屏幕上会提示您判断该图片出现过还是未出现过需要点击左边按钮还是右边按钮</p>
    <p>在整个过程中，都将保持同样的按键模式，请在保证准确率的前提下尽快按键反应</p>
	<p>如果无法判断或者不能确定请按“未出现”按键</p>
    <p>如果已经清楚，请点击开始按钮开始~ </p>
    <button class="btn btn-yes" id="start">开始</button>
    `,
    choices: "NO_KEYS",
    on_load: function () {
      document.getElementById("start").addEventListener("click", function () {
        jsPsych.finishTrial();
      });
    },
  };
  temp.push(welcome_trial);

  let lastNames = item.last.map((c) => c.name);

  let isEnd = false;
  for (let i = 0; i < item.list.length; i++) {
    isEnd = i == item.list.length - 1;
    // 空白
    const blank_trial = {
      type: jsPsychHtmlKeyboardResponse,
      stimulus: `
                    <div class="focus-container">
                        <div style="height: 28.8px;"></div>
                        <div class="fixation"> </div>
                    </div>
                `,
      choices: "NO_KEYS",
      trial_duration: blankDelayTime,
    };
    temp.push(blank_trial);

    // 注视点
    const focus_trial = {
      type: jsPsychHtmlKeyboardResponse,
      stimulus: `
                    <div class="focus-container">
                        <div style="height: 28.8px;"></div>
                        <div class="fixation">+</div>
                    </div>
                `,
      choices: "NO_KEYS",
      trial_duration: focusDelayTime,
    };
    temp.push(focus_trial);

    // 显示图片
    const image_trial = {
      type: jsPsychHtmlKeyboardResponse,
      stimulus: `<div >
            <div class="instruction">是否出现过</div>
            <img src="${item.list[i].link}">
           <div class="progress-container">
              <div class="progress-bar" style="width:${
                (i / item.list.length) * 100
              }%"></div>
            </div>
        </div>`,
      trial_duration: !isEnd ? block2ShowTime : (block2ShowTime + block2MaskTime),
      response_ends_trial: false,
    };
    temp.push(image_trial);
    // 显示mask图片4秒
    const maskImage = mask[jsPsych.randomization.randomInt(0, mask.length - 1)];
    if (!isEnd) {
      const mask_trial = {
        type: jsPsychHtmlKeyboardResponse,
        stimulus: `<div >
              <div class="instruction">是否出现过</div>
              <img src="${maskImage}">
             <div class="progress-container">
                <div class="progress-bar" style="width:${
                  (i / item.list.length) * 100
                }%"></div>
              </div>
          </div>`,
        trial_duration: block2MaskTime,
        response_ends_trial: false,
      };
      temp.push(mask_trial);
    }

    // 显示注视点和选择按钮
    const response_trial = {
      type: jsPsychHtmlKeyboardResponse,
      stimulus: `
                    <div class="focus-container">
                        <div class="instruction" id="instruction">是否出现过？</div>
                        <div class="fixation">+</div>
                        <div class="buttons">
                            ${
                              rand
                                ? [
                                    '<button class="btn btn-yes" id="btn-yes">出现过(F)</button>',
                                    '<button class="btn btn-no" id="btn-no">未出现(J)</button>',
                                  ].join("")
                                : [
                                    '<button class="btn btn-no" id="btn-no">未出现(F)</button>',
                                    '<button class="btn btn-yes" id="btn-yes">出现过(J)</button>',
                                  ].join("")
                            }
                        </div>
                    </div>
                `,
      choices: ["f", "j"],
      // extensions: isEnd ? [{ type: Naodao }] : [],
      data: {
        block: 2,
        trial: i + 1,
        实际出现过: Number(lastNames.includes(item.list[i].name)) || "",
        测谎图: Number(item.list[i].name.includes("测谎图")) || "",
        图片: item.list[i].name,
      },
      trial_duration: null,
      on_start: function (d) {
        d.data.startTime = new Date().getTime();
      },
      on_load: function () {
        clearTimeout(timer);
        timer = setTimeout(() => {
          let x = document.getElementById("btn-yes");
          let y = document.getElementById("btn-no");
          x.classList.add("overTime");
          y.classList.add("overTime");
          x.innerText = "继续(F)";
          y.innerText = "继续(J)";

          let instruction = document.getElementById("instruction");
          instruction.innerText = `你未在${
            block2DelayTime / 1000
          }秒内作答，操作已超时。按F/J继续`;
        }, block2DelayTime);
        // 添加按钮事件
        document
          .getElementById("btn-yes")
          .addEventListener("click", function () {
            jsPsych.finishTrial({ response: "yes" });
          });
        document
          .getElementById("btn-no")
          .addEventListener("click", function () {
            jsPsych.finishTrial({ response: "no" });
          });
      },
      on_finish: function (data) {
        data.endTime = new Date().getTime();
        data.diffTime = (data.endTime - data.startTime) + (block2ShowTime + block2MaskTime);
        data.endTime = formatDate(data.endTime);
        data.startTime = formatDate(data.startTime);
        clearTimeout(timer);
        timer = null;
        // 记录键盘响应
        if (data.response === "f") {
          data.response = rand ? "yes" : "no";
        } else if (data.response === "j") {
          data.response = rand ? "no" : "yes";
        }
      },
    };
    temp.push(response_trial);
  }
  timeline.push({
    timeline: temp,
    conditional_function: function () {
      const ch1 = jsPsych.data
        .get()
        .filter({ block: 1, 图片: "测谎图" })
        .last(1)
        .values()[0];
      if (ch1 && ch1.response === "no") {
        return false;
      }
      return true;
    },
  });
}
