module.exports = {
    reporter: function (res) {
        const sgr = (c) => "\x1b"+"["+c+"m";
        const reset = sgr(0);
        const gray = (s) => sgr(2) + s + reset;
        const green = (s) => sgr(32) + s + reset;
        const cyan = (s) => sgr(36) + s + reset;

        for (const item of res) {
            const file = item.file;
            const err = item.error;
            console.log(file + ":" + green(err.line) + ":" + cyan(err.character) + ", " + err.reason);
        }
    }
};
