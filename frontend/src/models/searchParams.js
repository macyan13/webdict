export default class SearchParams {
    constructor(tagIds, langId, page, pageSize, sourcePart = '', targetPart = '') {
        this.tagId = tagIds
        this.langId = langId;
        this.page = page;
        this.pageSize = pageSize;
        this.sourcePart = sourcePart;
        this.targetPart = targetPart;
    }
}