/** 文本标注任务 */
export interface ItfTextMarkTask {
  /** 样本id */
  datasetId: string;
  /** 数据序列 */
  dataSequence: number[];
  /** 任务类型 */
  annotationType: string;
  /** 任务名称 */
  name: string;
  /** 负责人 */
  principal: string;
}
