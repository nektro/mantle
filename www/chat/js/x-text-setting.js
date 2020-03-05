"use strict";
//
import { create_element, dcTN, setDataBinding } from "./util.js";

//
customElements.define("x-text-setting", class TextSetting extends HTMLElement {
    constructor() {
        super();
    }
    connectedCallback() {
        const e = this.getAttribute("endpoint");
        const n = this.getAttribute("name");
        const t = Math.random().toString().replace(".","");
        const v = this.getAttribute("value")||"";
        const d = this.getAttribute("label")||"";
        const b = this.getAttribute("binding");
        this.appendChild(create_element("form", [["method","post"],["action",e]], [
            create_element("label", [["for","input_"+t]], [dcTN(d)]),
            create_element("div", null, [
                create_element("input", [["type","text"],["name",n],["id","input_"+t],["value",v]]),
                create_element("button", null, [
                    create_element("i", [["class","check icon"]])
                ]),
            ]),
        ]));
        this.children[0].addEventListener("submit", (ev) => {
            ev.preventDefault();
            const fd = new FormData();
            const iv = this.querySelector("input").value;
            fd.append("p_name", n);
            fd.append("p_value", iv);
            return fetch(e, { method: "post", body: fd, }).then((x) => x.json()).then((x) => {
                console.info(x);
                if (b === null) return;
                setDataBinding(b, iv);
            });
        });
    }
    static get observedAttributes() {
        return ["value"];
    }
    attributeChangedCallback(name, oV, nV) {
        if (oV === null) {
            return;
        }
        if (name === "value") {
            this.querySelector("input").value = nV;
        }
    }
});
