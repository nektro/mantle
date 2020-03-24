"use strict";
//
import { WSetting } from "./w-setting.js";
import { create_element, dcTN, setDataBinding } from "./../util.js";

//
customElements.define("x-text-setting", class TextSetting extends WSetting {
    constructor() {
        super();
    }
    connectedCallback() {
        const n = this.getAttribute("name");
        const t = Math.random().toString().replace(".","");
        const v = this.getAttribute("value")||"";
        const d = this.getAttribute("label")||"";
        const b = this.getAttribute("binding");
        this.appendChild(create_element("form", null, [
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
            const de = this.defaultEndpoint();
            const e = this.getAttribute("endpoint")||de;
            const f = this.getAttribute("fill")||"";
            const e2 = e.replace("%s", f);
            const fd = new FormData();
            const iv = this.querySelector("input").value;
            fd.append("p_name", n);
            fd.append("p_value", iv);
            return fetch(e2, { method: "put", body: fd, }).then((x) => x.json()).then(() => {
                if (n === "name") {
                    this.parentElement.parentElement.children[0].querySelector(`[data-uid="${f}"]`).textContent = iv;
                }
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
