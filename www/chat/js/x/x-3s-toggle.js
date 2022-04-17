"use strict";
//
import { WSetting } from "./w-setting.js";
import { create_element, dcTN } from "./../util.js";

//
// jshint -W098
customElements.define("x-3s-toggle", class extends WSetting {
    constructor() {
        super();
    }

    connectedCallback() {
        const n = this.getAttribute("name");
        const d = this.getAttribute("label") || "";
        const t = [
            Math.random().toString().replace(".", ""),
            Math.random().toString().replace(".", ""),
            Math.random().toString().replace(".", ""),
        ];
        this.appendChild(create_element("form", [], [
            create_element("label", [], [dcTN(d)]),
            create_element("div", [["clasS", "switch-toggle switch-candy"]], [
                create_element("input", [["id", `deny_${t[0]}`], ["name", n], ["type", "radio"]]),
                create_element("label", [["for", `deny_${t[0]}`]], [
                    create_element("i", [["class", "times icon"]]),
                ]),
                create_element("input", [["id", `ignore_${t[1]}`], ["name", n], ["type", "radio"]]),
                create_element("label", [["for", `ignore_${t[1]}`]], [
                    create_element("i", [["class", "minus icon"]]),
                ]),
                create_element("input", [["id", `allow_${t[2]}`], ["name", n], ["type", "radio"]]),
                create_element("label", [["for", `allow_${t[2]}`]], [
                    create_element("i", [["class", "check icon"]]),
                ]),
            ]),
        ]));
        this.children[0][n].forEach((v) => {
            v.addEventListener("change", () => {
                const de = this.defaultEndpoint();
                const e = this.getAttribute("endpoint") || de;
                const f = this.getAttribute("fill") || "";
                const e2 = e.replace("%s", f);
                const fd = new FormData();
                fd.append("p_name", n);
                fd.append("p_value", (this.querySelector("input:checked").indexOfMe() / 2).toString());
                return fetch(e2, { method: "put", body: fd, });
            });
        });
    }

    static get observedAttributes() {
        return ["value"];
    }

    attributeChangedCallback(name, oV, nV) {
        if (name === "value") {
            const b = parseInt(nV, 10);
            if (b < 0 || b > 2) return;
            this.querySelectorAll("input")[b].checked = true;
        }
    }
});
