export default class User {
    constructor(id, source, transcription, target, example, tag_ids, langId) {
        this.id = id;
        this.source = source;
        this.transcription = transcription;
        this.target = target;
        this.example = example;
        this.tag_ids = tag_ids;
        this.lang_id = langId
    }
}