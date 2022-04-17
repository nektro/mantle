"use strict";
//
import { WSetting } from "./w-setting.js";
import { create_element, dcTN } from "./../util.js";

//
customElements.define("x-button", class extends WSetting {
    constructor() {
        super();
    }

    connectedCallback() {
        this.appendChild(create_element("form", null, [
            create_element("button", null, [dcTN(this.getAttribute("text"))]),
        ]));
        this.querySelector("form").addEventListener("submit", async (ev) => {
            ev.preventDefault();
            const de = this.defaultEndpoint();
            const e = this.getAttribute("endpoint") || de;
            const m = this.getAttribute("method") || "post";
            const f = this.getAttribute("fill") || "";
            const e2 = e.replace("%s", f);
            if (m === "delete") {
                const result = await Swal.fire({
                    title: "Are you sure?",
                    showCancelButton: true,
                    confirmButtonColor: "#3085d6",
                    cancelButtonColor: "#d33",
                    confirmButtonText: "Yes, delete it!"
                });
                if (!result.value) { return; }
            }
            return fetch(e2, { method: m, body: new FormData() });
        });
    }
});
