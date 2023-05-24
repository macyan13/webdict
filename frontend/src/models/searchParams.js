export default class SearchParams {
    constructor(tagId, langId, page, pageSize) {
        this.tagId = tagId
        this.langId = langId;
        this.page = page;
        this.pageSize = pageSize;
    }
}