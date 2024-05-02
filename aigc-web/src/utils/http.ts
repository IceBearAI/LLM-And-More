import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosError, type AxiosResponse } from "axios";
import { ElLoading } from "element-plus";
import { get, merge } from "lodash";
import awaitTo from "await-to-js";
import { useUserStore } from "@/stores";
import { toast } from "vue3-toastify";
import { useAppStore } from "@/stores";
import { saveAs } from "file-saver";

let mapLoadingInstance = {};

/**
 * showSuccess: 是否显示成功信息
 */
type MyConfig = {
  /**
   * 接口调用异常，是否使用默认提示方案
   * 1. false 不使用默认方案，接口自行处理
   * 2. true ， 使用默认方案，显示时长 3秒
   * 3. number， 使用默认方案，显示时长为传入的 number
   */
  showError?: boolean | number;
  /**
   * 接口调用成功后，是否显示成功信息
   * 1. false 不提示
   * 2. true：提示 “操作成功”
   * 3. 传入字符串，提示内容变更为 传入的字符串
   */
  showSuccess?: boolean | string;
  /**
   * 接口调用过程中，是否显示 loading
   * 1. true 显示范围整个body
   * 2. false 不显示
   * 3. css 选择器
   * 4. html 元素
   * 5. btn#btnSubmit 按钮
   */
  showLoading?: boolean | string | HTMLElement;
  /** loading实例id，无需传入，程序计算获得 */
  loadingId?: string;
};
type DownloadByUrlParams = {
  /** 文件下载地址 */
  fileUrl: string;
  /** 文件名称 */
  fileName?: string;
  /** 文件后缀名 */
  suffixName?: string;
};
type MyAxiosRequestConfig = AxiosRequestConfig & MyConfig;

type MyAxiosResponse = AxiosResponse & MyConfig;

/** 退出登录并强制刷新页面（会重定向到登录页） */
function logout() {
  toast.error("会话超时，请重新登录");
  useUserStore().logout();
}

const hideLoading = (config = { showLoading: true } as any) => {
  if (/^btn/.test(config.showLoading)) {
    const appStore = useAppStore();
    appStore.setBtnLoading(false);
  }
  if (mapLoadingInstance[config.loadingId]) {
    mapLoadingInstance[config.loadingId].close();
    delete mapLoadingInstance[config.loadingId];
  }
};

/**
 *
 * @param {*} param0
 * config : {
 *   showError 是否显示异常信息
 * }
 * @returns
 */
const errorHandler = ({ message, config, code, data }) => {
  let messageText = message || "Error";
  if (/502/.test(messageText)) {
    //后台服务重启中
    messageText = "系统服务升级中，请稍后重试";
  }
  let { showError } = config;
  if (showError) {
    toast.error(messageText, {
      autoClose: isNaN(parseInt(showError)) ? 3000 : parseInt(showError)
    });
  }
  return Promise.reject({
    message: messageText,
    code,
    data,
    toString() {
      return messageText;
    }
  });
};

/**
 * 生成loading 实例，唯一Id
 * @returns
 */
const getLoadingId = () => {
  let uuid = Math.random().toString(36).substring(2);
  if (mapLoadingInstance[uuid]) {
    //已存在，重新生成
    return getLoadingId();
  } else {
    return uuid;
  }
};

/** 创建请求实例 */
function createService() {
  // 创建一个 axios 实例命名为 service
  const service = axios.create();
  // 请求拦截
  service.interceptors.request.use(
    config => {
      config.url = ((apiOrigin || "") + import.meta.env.VITE_APP_BASE_API + config.url).replace("/api/api", "/api");
      return config;
    },
    // 发送失败
    error => Promise.reject(error)
  );
  // 响应拦截（可根据具体业务作出相应的调整）
  service.interceptors.response.use(
    async response => {
      // apiData 是 api 返回的数据
      let config = response.config as MyAxiosResponse;

      hideLoading(config);

      const apiData = response.data;
      // 二进制数据则直接返回
      const responseType = response.request?.responseType;
      if ("blob" == responseType) {
        return response;
      } else if ("arraybuffer" == responseType) {
        return { apiData };
      } else if (["string"].includes(typeof apiData)) {
        //聊天，返回流格式
        return apiData;
      }

      // 这个 code 是和后端约定的业务 code
      const { code, message, status, data } = apiData;

      // 如果没有 code, 代表这不是项目后端开发的 api
      if (code === undefined) {
        toast.error("非本系统的接口");
        return Promise.reject(new Error("非本系统的接口"));
      }

      switch (code) {
        case 200:
          // 本系统采用 code === 0 来表示没有业务错误
          let t = useAppStore().methods.t;
          if (config.showSuccess) {
            let messageText = t("api.success");
            if (typeof config.showSuccess == "string") {
              messageText = config.showSuccess;
            }
            toast.success(messageText);
          }
          return apiData.data || {};
        case 501:
          // Token 过期时
          return logout();
        default:
          // 不是正确的 code
          return errorHandler({ message, config, code, data });
      }
    },
    error => {
      // status 是 HTTP 状态码
      let config = error.config;
      hideLoading(config);
      // const status = get(error, "response.status");
      // switch (status) {
      //   case 400:
      //     error.message = "请求错误";
      //     break;
      //   case 401:
      //     // Token 过期时
      //     logout();
      //     break;
      //   case 403:
      //     error.message = "拒绝访问";
      //     break;
      //   case 404:
      //     error.message = "请求地址出错";
      //     break;
      //   case 408:
      //     error.message = "请求超时";
      //     break;
      //   case 500:
      //     error.message = "服务器内部错误";
      //     break;
      //   case 501:
      //     error.message = "服务未实现";
      //     break;
      //   case 502:
      //     error.message = "网关错误";
      //     break;
      //   case 503:
      //     error.message = "服务不可用";
      //     break;
      //   case 504:
      //     error.message = "网关超时";
      //     break;
      //   case 505:
      //     error.message = "HTTP 版本不受支持";
      //     break;
      //   default:
      //     break;
      // }
      switch (error.code) {
        // case "ECONNABORTED":
        //   error.message = "网络存在波动，若有异常请重试";
        //   break;

        default:
          break;
      }
      return errorHandler({
        message: error.message,
        config,
        code: error.code,
        data: error.data
      });
    }
  );
  return service;
}

/** 创建请求方法 */
function createRequest(service: AxiosInstance) {
  return function <T>(config: MyAxiosRequestConfig): Promise<[AxiosError, T]> {
    let defaults: MyConfig = {
      showError: true, //默认显示异常
      showSuccess: false,
      showLoading: false
    };
    config = { ...defaults, ...config };
    const { token, tenantId } = useUserStore().userInfo;
    const defaultConfig = {
      headers: {
        // 携带 Token
        "X-Token": token,
        "X-Tenant-Id": tenantId,
        "Content-Type": "application/json"
      },
      timeout: config.url == "/files" ? 60 * 1000 * 10 : 60 * 1000, // 请求超时时间 1000=1s
      // baseURL: import.meta.env.VITE_APP_BASE_API,
      data: {}
    };
    // 将默认配置 defaultConfig 和传入的自定义配置 config 进行合并成为 mergeConfig
    const mergeConfig = merge(defaultConfig, config);
    let { showLoading } = config;
    if (showLoading) {
      let el;
      if (typeof showLoading == "boolean") {
        el = "body";
      } else if (showLoading instanceof Element) {
        el = showLoading;
      } else if (/^btn/.test(showLoading)) {
        const appStore = useAppStore();
        appStore.setBtnLoading(true, showLoading.replace(/^btn#/, ""));
        el = null;
      } else {
        //css选择器
        el = showLoading;
      }

      if (el) {
        let loadingId = getLoadingId();
        mergeConfig.loadingId = loadingId;
        mapLoadingInstance[loadingId] = ElLoading.service({
          target: el
        });
      }
    }
    return awaitTo(service(mergeConfig)) as any;
  };
}

/** 用于网络请求的实例 */
const service = createService();
/** 用于网络请求的方法 */
// export const http = createRequest(service);
const instance = createRequest(service);

export const http = {
  get<T = any>(options: MyAxiosRequestConfig) {
    let { data = {}, url } = options;
    const arr = [];
    Object.keys(data).forEach(item => {
      if (Array.isArray(data[item])) {
        //数组转为 arr=item1&item2&item3
        const url = data[item].map(_ => `${item}=${_}`).join("&");
        arr.push(url);
      } else {
        if ([null, undefined, ""].includes(data[item])) {
          //无效数据，跳过
          return;
        }
        arr.push(`${item}=${data[item]}`);
      }
    });
    if (arr.length) {
      if (/\?$/.test(url)) {
        url += arr.join("&");
      } else if (/\?.+/.test(url)) {
        url += "&" + arr.join("&");
      } else {
        url += "?" + arr.join("&");
      }
    }

    delete options.data; //axios get请求用到的参数字段是params
    return instance<T>({
      ...options,
      url,
      method: "get"
    });
  },
  post<T = any>(options: MyAxiosRequestConfig) {
    return instance<T>({
      ...options,
      method: "post"
    });
  },
  put<T = any>(options: MyAxiosRequestConfig) {
    return instance<T>({
      ...options,
      method: "put"
    });
  },
  delete<T = any>(options: MyAxiosRequestConfig) {
    let { data } = options;
    delete options.data; //axios get请求用到的参数字段是params
    return instance<T>({
      ...options,
      params: {
        ...data
      },
      method: "delete"
    });
  },
  /**
   * 文件上传
   * @param options
   * @param isMultiple  boolean 是否为多文件上传，默认false
   * @returns
   */
  upload<T = any>(options: MyAxiosRequestConfig, isMultiple = false) {
    let { data } = options;
    let formData = new FormData();
    for (let key in data) {
      let fieldValue = data[key];
      if (Array.isArray(fieldValue) && fieldValue[0] instanceof File) {
        if (isMultiple) {
          formData.append(key, data[key]);
        } else {
          //单文件上传
          formData.append(key, data[key][0]);
        }
      } else {
        formData.append(key, data[key]);
      }
    }

    return instance<T>({
      ...options,
      method: "post",
      data: formData,
      headers: {
        "content-type": "multipart/form-data"
      }
    });
  },
  async download<T>(options: MyAxiosRequestConfig) {
    let defaults = {
      showLoading: true,
      method: "get", //默认使用get方式
      timeout: Number.MAX_SAFE_INTEGER
    };
    options = {
      ...defaults,
      ...options,
      responseType: "blob"
    };
    const [err, res] = await this[options.method](options);
    if (res) {
      try {
        let blob, fileName;
        fileName = decodeURIComponent(res.headers["content-disposition"].split("=")[1].replace(/"/g, ""));
        blob = new Blob([res.data], {
          type: "application/zip"
        });
        saveAs(blob, fileName);
      } catch (e) {
        console.error("下载失败");
      }
      return [null, {}];
    } else {
      return [err];
    }
  },
  downloadByUrl: (options: DownloadByUrlParams) => {
    let { fileUrl, fileName, suffixName } = options;
    fetch(fileUrl)
      .then(resp => resp.blob())
      .then(blob => {
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement("a");
        a.style.display = "none";
        a.href = url;
        let downloadFileName: string = +new Date() + ""; //文件下载到本地以后的名称
        if (fileName) {
          downloadFileName = fileName;
        } else if (suffixName) {
          downloadFileName += "." + suffixName;
        }
        a.download = downloadFileName;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
      });
  }
};
