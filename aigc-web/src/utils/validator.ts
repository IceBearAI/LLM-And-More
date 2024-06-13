type typeValidator = {
  value: string;
  errorEmpty?: string;
  required?: boolean;
  errorValid?: string;
};

type typeValidatorPrivate = {
  value: string;
  exp: RegExp;
  required?: boolean;
  errorEmpty?: string;
  errorValid: string;
};

const handlerValidator = ({ value, exp, required, errorEmpty, errorValid }: typeValidatorPrivate) => {
  if (required) {
    //必填
    if (value) {
      // let reg = new RegExp(exp);
      if (exp.test(value)) {
        return true;
      } else {
        return errorValid;
      }
    } else {
      //空值，返回 errorEmpty
      return errorEmpty;
    }
  } else {
    //非必填
    if (value) {
      if (exp.test(value)) {
        return true;
      } else {
        return errorValid;
      }
    } else {
      //空值
      return true;
    }
  }
};

export const validator = {
  isModelName({
    value,
    errorEmpty = "请输入模型名称",
    required = true,
    errorValid = "只允许字母、数字、“-” 、“.” 和 “:” "
  }: typeValidator) {
    return handlerValidator({
      value,
      exp: /^[a-zA-Z0-9-.:]+$/,
      required,
      errorEmpty,
      errorValid
    });
  },
  isEmail({ value, errorEmpty = "请输入邮箱", required = false, errorValid = "请输入有效的邮箱" }: typeValidator) {
    return handlerValidator({
      value,
      exp: /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/,
      required,
      errorEmpty,
      errorValid
    });
  },
  isPassword({
    value,
    errorEmpty = "请输入密码",
    required = false,
    errorValid = "不符合密码要求格式：大写英文字母+小写英文字母+数字+特殊字符，8-20位"
  }: typeValidator) {
    return handlerValidator({
      value,
      exp: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,20}$/,
      required,
      errorEmpty,
      errorValid
    });
  },
  //中文、空格、标点符号
  isCNOrSymbol({ value, errorEmpty = "请输入", required = false, errorValid = "请输入中文" }: typeValidator) {
    return handlerValidator({
      value,
      exp: /^[\u4E00-\u9FA5\uF900-\uFA2D\s,.!();:?，。！（）；：？、]+$/,
      required,
      errorEmpty,
      errorValid
    });
  },
  // 中文、数字、字母、-、_
  isName({ value, errorEmpty = "请输入", required = false, errorValid = "请输入正确格式" }: typeValidator) {
    return handlerValidator({
      value,
      exp: /^[A-Za-z0-9-_\u4e00-\u9fa5]+$/, ///^[A-Za-z0-9-_\u4e00-\u9fa5]{4,30}$/
      required,
      errorEmpty,
      errorValid
    });
  },
  validNumberInput(value, min, max, errorMessage, reg = false) {
    if (value) {
      if (value < min) {
        return `下限 ${min}`;
      } else if (value > max) {
        return `上限 ${max}`;
      } else if (reg && /^\+?[1-9][0-9]*$/.test(value) != true) {
        return "请输入正整数";
      } else {
        return true;
      }
    } else {
      if (errorMessage) {
        return errorMessage;
      }
    }
    return true;
  },
  isJson(str) {
    if (typeof str == "string") {
      try {
        var obj = JSON.parse(str);
        // 等于这个条件说明就是JSON字符串 会返回true
        if (typeof obj == "object" && obj) {
          return true;
        } else {
          //不是就返回false
          return false;
        }
      } catch (e) {
        return false;
      }
    }
    return false;
  }
};
