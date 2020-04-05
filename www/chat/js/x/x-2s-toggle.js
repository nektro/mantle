"use strict";
//
import { WSetting } from "./w-setting.js";
import { create_element, dcTN } from "./../util.js";

//
// jshint -W098
customElements.define("x-2s-toggle", class extends WSetting {
    constructor() {
        super();
    }
    connectedCallback() {
        const n = this.getAttribute("name");
        const d = this.getAttribute("label")||"";
        this.appendChild(create_element("div", null, [dcTN(d)]));
        this.appendChild(create_element("label", null, [
            create_element("input", [["type","checkbox"]]),
            create_element("span")
        ]));
        this.children[1].children[0].addEventListener("change", (ev) => {
            const de = this.defaultEndpoint();
            const e = this.getAttribute("endpoint")||de;
            const f = this.getAttribute("fill")||"";
            const e2 = e.replace("%s", f);
            const fd = new FormData();
            fd.append("p_name", n);
            fd.append("p_value", ev.target.checked ? "1" : "0");
            return fetch(e2, { method: "put", body: fd, });
        });
    }
    static get observedAttributes() {
        return ["value"];
    }
    attributeChangedCallback(name, oV, nV) {
        if (name === "value") {
            const b = nV === "true" || nV === "1";
            const v = this.hasAttribute("inverted");
            this.children[1].children[0].checked = v ? !b : b;
        }
    }
});
