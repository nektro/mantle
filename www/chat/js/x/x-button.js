"use strict";
//
import { WSetting } from "./w-setting.js";
import { create_element, dcTN } from "./../util.js";

//
customElements.define("x-button", class XButton extends WSetting {
    constructor() {
        super();
    }
    connectedCallback() {
        this.appendChild(create_element("form", null, [
            create_element("button", null, [dcTN(this.getAttribute("text"))]),
        ]));
        this.querySelector("form").addEventListener("submit", (ev) => {
            ev.preventDefault();
            const de = this.defaultEndpoint();
            const e = this.getAttribute("endpoint")||de;
            const m = this.getAttribute("method")||"post";
            const f = this.getAttribute("fill")||"";
            const e2 = e.replace("%s", f);
            return fetch(e2, { method:m, body:new FormData() });
        });
    }
});
