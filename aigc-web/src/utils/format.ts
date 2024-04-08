import moment from "moment";
import Big from "big.js";
import { useAppStore } from "@/stores";
import { EnumLanguage } from "@/types/apps/System.ts";

export const format = {
  // 手机号脱敏，星号替换
  blurMobile(inValue) {
    return ("" + inValue).replace(/^(\d{3}).+(\d{4})$/, "$1****$2");
  },

  /**
   * 千分号字符串
   *  eg 123 -> 123  、  1234 -> 1,234      、     1234567 -> 1,234,567
   * @param inValue
   * @returns
   */
  commaString(inValue) {
    return inValue.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
  },

  /**
   * 日期格式化
   * @param {*} rawDate 日期初始值
   * @param {*} formatStr 日期格式标准  "YYYY-MM-DD HH:mm:ss"
   * @returns
   */
  dateFormat(rawDate, formatStr = "YYYY-MM-DD") {
    if (rawDate) {
      const appStore = useAppStore();
      let { localLanguage } = appStore;
      let ret = "";
      if (localLanguage == EnumLanguage.English) {
        //英国时间，比中国慢8个小时
        ret = moment(rawDate).subtract(8, "hour").format(formatStr);
      } else if (localLanguage == EnumLanguage.Espanyol) {
        //西班牙时间，比中国慢7个小时
        ret = moment(rawDate).subtract(7, "hour").format(formatStr);
      } else if (localLanguage == EnumLanguage.Pilipino) {
        //菲律宾时间，跟中国时间等同
        ret = moment(rawDate).format(formatStr);
      } else {
        //默认使用中国时间
        ret = moment(rawDate).format(formatStr);
      }
      if (ret == "Invalid date") {
        return "";
      } else {
        return ret;
      }
    } else {
      return "";
    }
  },
  getFileSize(size) {
    if (!size) return "";
    const num = 1024.0;
    if (size < num) return size + "B";
    if (size < Math.pow(num, 2)) return (size / num).toFixed(2) + "K"; //kb
    if (size < Math.pow(num, 3)) return (size / Math.pow(num, 2)).toFixed(2) + "M"; //M
    if (size < Math.pow(num, 4)) return (size / Math.pow(num, 3)).toFixed(2) + "G"; //G
    return (size / Math.pow(num, 4)).toFixed(2) + "T"; //T
  },
  toPercent(value, point = 0) {
    if (value === undefined || value === "") {
      return "";
    } else {
      if (value == 1) {
        return "100%";
      } else {
        return (value * 100).toFixed(point) + "%";
      }
    }
  },
  /**
   * number转Scientific
   * @param {String} number 数字
   * @param {Number} power 到多少位才转换
   * @returns 科学记数法
   */
  toScientfic(number, power = 4) {
    if (number && typeof number === "string") {
      if (number.toString().includes("e")) return number;
      const value = Number(number);
      if (value.toString().includes("e")) return value.toString();
      const p = Math.floor(Math.log(Math.abs(value)) / Math.LN10);
      if (Math.abs(p) < power) return value.toString();
      const bigNumber = new Big(value);
      return bigNumber.toExponential();
    } else return "";
  },
  dateFromNow(date) {
    return moment(date).fromNow();
  }
};
