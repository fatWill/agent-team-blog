let resource = {
  mask: [
    "./images/mask/mask1.png",
    "./images/mask/mask2.png",
    "./images/mask/mask3.png",
    "./images/mask/mask4.png",
  ],
  test: [
    {
      link: "./images/test/1402.png",
      name: "1402",
    },
    {
      link: "./images/test/1804.png",
      name: "1804",
    },
    {
      link: "./images/test/2203.png",
      name: "2203",
    },
  ],
  main: [
    {
      link: "./images/main/0307.png",
      name: "0307",
    },
    {
      link: "./images/main/0308.png",
      name: "0308",
    },
    {
      link: "./images/main/0407.png",
      name: "0407",
    },
    {
      link: "./images/main/0408.png",
      name: "0408",
    },
    {
      link: "./images/main/62701.jpg",
      name: "62701",
    },
    {
      link: "./images/main/62702.jpg",
      name: "62702",
    },
    {
      link: "./images/main/62703.jpg",
      name: "62703",
    },
    {
      link: "./images/main/62704.jpg",
      name: "62704",
    },
    {
      link: "./images/main/62705.jpg",
      name: "62705",
    },
    {
      link: "./images/main/62706.jpg",
      name: "62706",
    },
    {
      link: "./images/main/62707.jpg",
      name: "62707",
    },
    {
      link: "./images/main/62801.jpeg",
      name: "62801",
    },
    {
      link: "./images/main/62802.png",
      name: "62802",
    },
    {
      link: "./images/main/62803.png",
      name: "62803",
    },
    {
      link: "./images/main/62806.png",
      name: "62806",
    },
    {
      link: "./images/main/62807.png",
      name: "62807",
    },
  ],
  lib: [
    {
      link: "./images/lib/0202.png",
      name: "0202",
    },
    {
      link: "./images/lib/0205.png",
      name: "0205",
    },
    {
      link: "./images/lib/0305.png",
      name: "0305",
    },
    {
      link: "./images/lib/0404.png",
      name: "0404",
    },
    {
      link: "./images/lib/0602.png",
      name: "0602",
    },
    {
      link: "./images/lib/0604.png",
      name: "0604",
    },
    {
      link: "./images/lib/1702.png",
      name: "1702",
    },
    {
      link: "./images/lib/1704.png",
      name: "1704",
    },
    {
      link: "./images/lib/1901.png",
      name: "1901",
    },
    {
      link: "./images/lib/1902.png",
      name: "1902",
    },
    {
      link: "./images/lib/2405.png",
      name: "2405",
    },
    {
      link: "./images/lib/2406.png",
      name: "2406",
    },
    {
      link: "./images/lib/2601.png",
      name: "2601",
    },
    {
      link: "./images/lib/2604.png",
      name: "2604",
    },
    {
      link: "./images/lib/3001.png",
      name: "3001",
    },
    {
      link: "./images/lib/3002.png",
      name: "3002",
    },
    {
      link: "./images/lib/3101.png",
      name: "3101",
    },
    {
      link: "./images/lib/3102.png",
      name: "3102",
    },
    {
      link: "./images/lib/4001.png",
      name: "4001",
    },
    {
      link: "./images/lib/4002.png",
      name: "4002",
    },
    {
      link: "./images/lib/40901.png",
      name: "40901",
    },
    {
      link: "./images/lib/40902.png",
      name: "40902",
    },
    {
      link: "./images/lib/40903.png",
      name: "40903",
    },
    {
      link: "./images/lib/4102.png",
      name: "4102",
    },
    {
      link: "./images/lib/4103.png",
      name: "4103",
    },
    {
      link: "./images/lib/4701.png",
      name: "4701",
    },
    {
      link: "./images/lib/4702.png",
      name: "4702",
    },
    {
      link: "./images/lib/4901.png",
      name: "4901",
    },
    {
      link: "./images/lib/4903.png",
      name: "4903",
    },
    {
      link: "./images/lib/5801.png",
      name: "5801",
    },
    {
      link: "./images/lib/5903.png",
      name: "5903",
    },
    {
      link: "./images/lib/6002.png",
      name: "6002",
    },
    {
      link: "./images/lib/6003.png",
      name: "6003",
    },
    {
      link: "./images/lib/62001.jpg",
      name: "62001",
    },
    {
      link: "./images/lib/62003.jpg",
      name: "62003",
    },
    {
      link: "./images/lib/6201.png",
      name: "6201",
    },
    {
      link: "./images/lib/6203.png",
      name: "6203",
    },
    {
      link: "./images/lib/62901.png",
      name: "62901",
    },
    {
      link: "./images/lib/62902.png",
      name: "62902",
    },
    {
      link: "./images/lib/62903.png",
      name: "62903",
    },
    {
      link: "./images/lib/6401.png",
      name: "6401",
    },
    {
      link: "./images/lib/6403.png",
      name: "6403",
    },
    {
      link: "./images/lib/9901.png",
      name: "9901",
    },
  ],
  compet: [],
};

let getResource = (jsPsych) => {
  // -------------------------回合1的素材-------------------------
  let block1MainResource = jsPsych.randomization
    .shuffle(resource.main)
    .slice(0, 12);

  let block1LibResource = jsPsych.randomization
    .shuffle(resource.lib)
    .slice(0, 20);

  let block1Resource = [
    ...jsPsych.randomization.shuffle(resource.test),
    ...jsPsych.randomization.shuffle([
      ...block1MainResource,
	  ...block1LibResource,
	]),
  ];

  // -------------------------回合2的素材-------------------------
  // 出现过的素材
  let block2OldMainResource = jsPsych.randomization
    .shuffle(block1MainResource)
    .slice(0, 6);
	
  let block2OldLibResource = jsPsych.randomization
    .shuffle(block1LibResource)
    .slice(0, 10);
  
  let block2OldResource = [
    ...block2OldMainResource,
	...block2OldLibResource,
  ];

  // 未出现的素材
  let block2OldMainNames = block1MainResource.map((c) => c.name);
  let block2NewMainResource = jsPsych.randomization.shuffle(
    resource.main.filter((c) => !block2OldMainNames.includes(c.name))
  )
  .slice(0,4);  
  
  let block2OldLibNames = block1LibResource.map((c) => c.name);
  let block2NewLibResource = jsPsych.randomization.shuffle(
    resource.lib.filter((c) => !block2OldLibNames.includes(c.name))
  )
  .slice(0,12);
  
  let block2NewResource = [
    ... block2NewMainResource,
	... block2NewLibResource,
  ];

  let block2Resource = [
    ...jsPsych.randomization.shuffle(resource.test),
    ...jsPsych.randomization.shuffle([
      ...block2OldResource,
      ...block2NewResource,
    ]),
  ];

  return {
    block1: {
      list: block1Resource.concat({
        link: "./images/xh.webp",
        name: "测谎图",
      }),
      last: [],
    },
    block2: {
      list: block2Resource.concat({
        link: "./images/cx.webp",
        name: "测谎图",
      }),
      last: [...block2OldResource],
    },
  };
};
