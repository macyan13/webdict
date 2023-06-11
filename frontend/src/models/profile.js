export default class Profile {
    constructor(id, name, email, currentPassword, newPassword, defaultLangId) {
        this.id = id;
        this.name = name;
        this.email = email;
        this.currentPassword = currentPassword;
        this.newPassword = newPassword;
        this.defaultLangId = defaultLangId;
    }
}
