"use strict";
//
import { create_element, dcTN } from "./../util.js";

//
customElements.define("x-button", class XButton extends HTMLElement {
    constructor() {
        super();
    }
    connectedCallback() {
        this.appendChild(create_element("form", null, [
            create_element("button", null, [dcTN(this.getAttribute("text"))]),
        ]));
        this.querySelector("form").addEventListener("submit", (ev) => {
            ev.preventDefault();
            const e = this.getAttribute("endpoint");
            const m = this.getAttribute("method")||"post";
            const f = this.getAttribute("fill")||"";
            const e2 = e.replace("%s", f);
            return fetch(e2, { method:m, body:new FormData() });
        });
    }
});
