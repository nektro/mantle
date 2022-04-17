"use strict";
//

//
export class WSetting extends HTMLElement {
    constructor() {
        super();
    }

    defaultEndpoint() {
        return `./../api/${this.parentElement.parentElement.parentElement.getAttribute("data-s-section")}/%s`;
    }
}
