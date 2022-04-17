"use strict";
//
import { WSetting } from "./w-setting.js";
import { create_element, dcTN } from "./../util.js";
import * as ui from "./../ui.js";

//
// jshint -W098
customElements.define("x-2s-toggle", class extends WSetting {
    constructor() {
        super();
    }

    connectedCallback() {
        const ln = this.getAttribute("local-name");
        const n = this.getAttribute("name");
        const d = this.getAttribute("label") || "";
        this.appendChild(create_element("div", null, [dcTN(d)]));
        this.appendChild(create_element("label", null, [
            create_element("input", [["type", "checkbox"]]),
            create_element("span")
        ]));
        if (ln !== null) {
            this.setAttribute("value", localStorage.getItem(ln));
            this.children[1].children[0].addEventListener("change", () => {
                ui.toggleHandlers.get(ln)(this._value);
            });
            return;
        }
        this.children[1].children[0].addEventListener("change", () => {
            const de = this.defaultEndpoint();
            const e = this.getAttribute("endpoint") || de;
            const f = this.getAttribute("fill") || "";
            const e2 = e.replace("%s", f);
            const fd = new FormData();
            fd.append("p_name", n);
            fd.append("p_value", this._value);
            return fetch(e2, { method: "put", body: fd, });
        });
    }

    get _value() {
        return this.children[1].children[0].checked ? "1" : "0";
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
