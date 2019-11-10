/**
 * @see https://github.com/nektro/mantle
 * @author Meghan Denny <https://nektro.net>
 */
//
import { el_1, el_2, el_3, getUserFromUUID } from "./util.js";

//
(async function() {
    //
    await fetch("/api/about").then(x => x.json()).then(x => {
        el_2.innerText = x.name;
    });

    //
    await fetch("/api/users/@me").then(x => x.json()).then(x => {
        if (x.success === false) {
            location.assign("../");
        }
        else {
            const u = x.message.me;
            const n = u.nickname || u.name;
            el_3.children[0].textContent = `@${n}`;

            const p = x.message.perms;
            for (const key in p) {
                document.querySelectorAll(`[data-requires^="${key}"]`).forEach((el) => {
                    el.removeAttribute("hidden");
                });
            }
        }
    });

    //
    await fetch("/api/channels/@me").then(x => x.json()).then(x => {
        for (const item of x.message) {
            ui.addChannel(item.uuid, item.name);
        }
        ui.setActiveChannel(x.message[0].uuid);

        const el2 = document.getElementById("channel-name");
        el2.children[0].textContent = x.message[0].name;
        el2.children[1].textContent = x.message[0].description;

        el_1.querySelector("button").addEventListener("click", async (e) => {
            const {value: name} = await Swal({
                title: "Enter the new channel's name",
                input: "text",
                showCancelButton: true,
                inputValidator: (value) => {
                    return !value && "You need to write something!"
                },
            });
            if (name !== undefined) {
                const fd = new URLSearchParams();
                fd.append("name", name);
                return fetch("/api/channels/create", {
                    method: "post",
                    credentials: "include",
                    headers: {
                        "content-type": "application/x-www-form-urlencoded",
                    },
                    body: fd,
                });
            }
        });
    });
})();
