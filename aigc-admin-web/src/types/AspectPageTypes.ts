export interface ItfAspectPageState {
  /** keep-alive组件，记录组件离开时的 scrollTop 信息 */
  scrollTop: Record<string, any>;
  methods: {
    /** 刷新列表页面，需列表页自己实现，其他页按需调用 */
    refreshListPage(): void;
  };
}
