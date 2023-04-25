export default class User {
    constructor(source, transcription, target, example, tag_ids) {
        this.source = source;
        this.transcription = transcription;
        this.target = target;
        this.example = example;
        this.tag_ids = tag_ids;
    }
}