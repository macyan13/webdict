export default function authHeader() {
    let user = JSON.parse(localStorage.getItem('user'));

    if (user && user.accessToken && user.type) {
        return { Authorization: user.type + ' ' + user.accessToken };
    } else {
        return {};
    }
}