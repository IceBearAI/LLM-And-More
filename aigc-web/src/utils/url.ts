export const url = {
  /**
   * 新窗口打开页面
   * @param pageUrl
   */
  onNewPage(pageUrl: string) {
    if (/^http/.test(pageUrl)) {
      //绝对路径
      window.open(pageUrl);
    } else {
      //相对路径，本项目页面
      if (location.hash) {
        //router history 设置为hash规则
        if (/^[^#]/.test(pageUrl)) {
          //pageUrl 未以 '#' 开头
          pageUrl = "#" + pageUrl;
        }
      }
      window.open(pageUrl);
    }
  }
};
