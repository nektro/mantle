"use strict";
//

//
customElements.define("x-message", class extends HTMLElement {
    constructor() {
        super();
    }
    get time() {
        return new moment(this.querySelector(".ts").getAttribute("title"), moment.defaultFormat);
    }
    connectedCallback() {
        this._uid = this.getAttribute("uuid");
        this._author = this.getAttribute("author");
    }
});
