export const file = {
  asyncImportJavaScript(u, isNew = true) {
    return new Promise((resolve, reject) => {
      var d = document,
        t = "script",
        o = d.createElement(t) as any,
        s = d.getElementsByTagName(t)[0];
      o.async = true;
      let src = u;
      if (!/^http/.test(src)) {
        //非绝对地址
        src = `//${u}`;
      }
      if (isNew) {
        // src += `?t=${+new Date()}`;
        src += `?t=${appVersion.replace(/[\s-:]+/g, "")}`;
      }
      o.src = src;

      o.addEventListener(
        "load",
        function (e) {
          resolve(e);
        },
        false
      );
      o.addEventListener(
        "error",
        function (e) {
          console.error("文件异步加载失败:" + u);
          reject(e);
        },
        false
      );
      s.parentNode.insertBefore(o, s);
    });
  }
};
