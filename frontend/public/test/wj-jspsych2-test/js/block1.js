/**
 * 注入回合1
 * @param {*} timeline
 * @param {*} rand
 * @param {
 *  list: any[],
 *  last: any[],
 * } item
 */
function injectBlock1(timeline, rand, item) {
  const welcome_trial = {
    type: jsPsychHtmlKeyboardResponse,
    stimulus: `
    <p>玩家您好，接下来在屏幕中会出现一张美术素材，请判断您是否喜欢该美术素材（如整体风格、细节设计等）</p>
    <p>屏幕上会提示您${
      rand ? "喜欢、不喜欢" : "不喜欢、喜欢"
    }该图片需要点击左边按钮还是右边按钮</p>
    <p>在整个过程中，都将保持同样的按键模式，请按照您对该图片的第一感觉来进行判断<br><br>如果已经清楚，请点击开始按钮开始~</p>
    <button class="btn btn-yes" id="start">开始</button>
    `,
    choices: "NO_KEYS",
    on_load: function () {
      document.getElementById("start").addEventListener("click", function () {
        jsPsych.finishTrial();
      });
    },
  };
  timeline.push(welcome_trial);

  for (let i = 0; i < item.list.length; i++) {
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
    timeline.push(blank_trial);

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
    timeline.push(focus_trial);

    // 显示图片
    const image_trial = {
      type: jsPsychHtmlKeyboardResponse,
      stimulus: `<div>
            <div class="instruction" id="instruction">是否喜欢</div>
            <img src="${item.list[i].link}">
           <div class="progress-container">
              <div class="progress-bar" style="width:${
                (i / item.list.length) * 100
              }%"></div>
            </div>
            <div class="buttons">
                ${
                  rand
                    ? [
                        '<button class="btn btn-yes" id="btn-yes">喜欢(F)</button>',
                        '<button class="btn btn-no" id="btn-no">不喜欢(J)</button>',
                      ].join("")
                    : [
                        '<button class="btn btn-no" id="btn-no">不喜欢(F)</button>',
                        '<button class="btn btn-yes" id="btn-yes">喜欢(J)</button>',
                      ].join("")
                }
            </div>
        </div>`,
      choices: ["f", "j"],
      trial_duration: null,
      data:{
        block:1,
        trial:(i+1),
        测谎图: Number(item.list[i].name.includes("测谎图")) || "",
        '图片':item.list[i].name
      },
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
            block1DelayTime / 1000
          }秒内作答，操作已超时。按F/J继续`;
        }, block1DelayTime);
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
        data.diffTime = data.endTime - data.startTime;
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
    timeline.push(image_trial);
  }
}
