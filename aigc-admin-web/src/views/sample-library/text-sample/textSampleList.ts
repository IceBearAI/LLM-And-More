export interface ItfModel {
  /** 主键，用于接口请求 */
  id: string;
  /** 名称 */
  modelName: string;
  /** 备注 */
  remark: string;
  /** 样本ID */
  sampleFileId: string;
}

type a = Partial<ItfModel>;
