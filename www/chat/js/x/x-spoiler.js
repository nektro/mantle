"use strict";
//

//
customElements.define("x-spoiler", class extends HTMLElement {
    constructor() {
        super();
    }
    connectedCallback() {
        this.addEventListener("click", () => {
            this.setAttribute("class", "show");
        });
    }
});
