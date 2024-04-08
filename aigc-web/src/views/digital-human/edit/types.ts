/** 数字人 */
export interface ItfDigitalHuman {
  /**  主键，用于接口提交 */
  name: string;
  /** 用户名，用于界面展示 */
  cname: string;
  /** 头像 */
  cover: string;
  /** 视频地址 */
  video: string;
  /** 类型 image图片、video视频 */
  mediaType: string;
}

/**  发声人 */
export interface ItfSpeaker {
  /**  主键，用于接口提交 */
  speakName: string;
  /** 用户名，用于界面展示 */
  speakCname: string;
  /** 头像 */
  headImg: string;
  /** 音频频地址 */
  speakDemo: string;
  /** 性别 */
  gender: number | "";
  /** 年龄段 */
  ageGroup: number | "";
  /** 用户展示信息，根据用户基本信息，计算而得 */
  subTitle: string;
}

export interface ItfProvideState {
  selectedDigitalHuman: ItfDigitalHuman;
  selectedSpeaker: ItfSpeaker;
  listSpeaker: ItfSpeaker[];
  [propName: string]: any;
}
