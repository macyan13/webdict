export default class SearchParams {
    constructor(tagIds, langId, page, pageSize) {
        this.tagId = tagIds
        this.langId = langId;
        this.page = page;
        this.pageSize = pageSize;
    }
}