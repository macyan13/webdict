export default class Profile {
    constructor(id, name, email, currentPassword, newPassword, defaultLangId, listOptions) {
        this.id = id;
        this.name = name;
        this.email = email;
        this.currentPassword = currentPassword;
        this.newPassword = newPassword;
        this.defaultLangId = defaultLangId;
        this.listOptions = listOptions;
    }
}

export class ListOptions {
    constructor(hideTranscription) {
        this.hideTranscription = hideTranscription;
    }
}
