/**
 * @see https://github.com/nektro/mantle
 * @author Meghan Denny <https://nektro.net>
 */
//
import { el_1, el_2, el_3, getUserFromUUID, output, messageCache } from "./util.js";
import * as ui from "./ui.js";

//
let me = null;

//
(async function() {
    //
    await fetch("./../api/about").then(x => x.json()).then(x => {
        console.info(x);
        el_2.innerText = x.message.name;
    });

    //
    await fetch("./../api/users/@me").then(x => x.json()).then(x => {
        console.info(x);
        if (x.success === false) {
            location.assign("../");
        }
        else {
            console.info(x.message);
            me = x.message.me;
            const n = me.nickname || me.name;
            el_3.children[0].textContent = `@${n}`;

            const p = x.message.perms;
            for (const key in p) {
                if (!p[key]) {
                    document.querySelectorAll(`[data-requires^="${key}"]`).forEach((el) => {
                        el.setAttribute("hidden", "");
                    });
                }
            }
        }
    });

    //
    await fetch("./../api/channels/@me").then(x => x.json()).then(async function(x) {
        console.info(x);
        for (const item of x.message) {
            console.info(item);
            await ui.addChannel(item);
        }
        await ui.setActiveChannel(x.message[0].uuid);

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
                return fetch("./../api/channels/create", {
                    method: "post",
                    credentials: "include",
                    headers: {
                        "content-type": "application/x-www-form-urlencoded",
                    },
                    body: fd,
                });
            }
        });
        el_1.querySelector("ol").addEventListener("click", (ev) => {
            ui.setActiveChannel(ev.target.dataset.uuid);
        });
        output.addEventListener("scroll", async function(e) {
            if (output.children.length === 0) { return }
            if (e.target.scrollTop !== 0) { return }
            if (output.classList.contains("loading")) { return }
            if (output.classList.contains("loading-done")) { return }
            //
            output.classList.add("loading");
            const fc = output.children[0];
            const lstm = output.children[0].dataset.msgUid;
            const chuid = ui.volatile.activeChannel.dataset.uuid;
            await fetch(`./../api/channels/${chuid}/messages?after=${lstm}`).then(x=>x.json()).then(async function(x) {
                if (x.message.length <= 1) {
                    output.classList.add("loading-done")
                    return
                }
                for (let i = 1; i < x.message.length; i++) {
                    const item = x.message[i];
                    const time = new Date(item.time.replace(" ","T")+"Z").toLocaleString();
                    output.prepend(ui.createMessage(await getUserFromUUID(item.author), {...item, time:time}))
                    messageCache.get(chuid).unshift(item);
                }
                output.scrollTop = fc.offsetTop-60;
            })
            output.classList.remove("loading");
        })
    });

    await fetch("./../api/users/online").then(x => x.json()).then(x => {
        console.info(x);
        for (const item of x.message) {
            ui.setMemberOnline(item);
        }
    });

    //
    const input = document.getElementById("input").children[0];
    const socket = new WebSocket(`ws${location.protocol.substring(4)}//${location.host}/ws`);

    socket.addEventListener("open", function() {
        el_2.classList.add("online");
        ui.addMessage(null, {name:"Connection Status"}, {body:"Active"}, false);
    });
    socket.addEventListener("close", function() {
        el_2.classList.remove("online");
    });
    socket.addEventListener("message", async function(e) {
        const d = JSON.parse(e.data);
        switch (d.type) {
            case "pong": {
                // do nothing, keep connection alive
                console.debug("pong")
                break;
            }
            case "message": {
                const u = await getUserFromUUID(d.message.author);
                ui.addMessage(d.in, u, d.message, true, Date.parse(d.at));
                break;
            }
            case "new-channel": {
                ui.addChannel(d);
                break;
            }
            case "user-connect": {
                ui.setMemberOnline(d.user);
                break;
            }
            case "user-disconnect": {
                ui.setMemberOffline(d.user);
                if (d.user === me.uuid) socket.close();
                break;
            }
            default: {
                console.log(d);
            }
        }
    });
    setInterval(function() {
        if (el_2.classList.contains("online")) {
            socket.send(JSON.stringify({
                type: "ping",
            }))
        }
    }, 30*1000)

    input.addEventListener("keydown", function(e) {
        if (e.key === "Enter") {
            socket.send(JSON.stringify({
                type: "message",
                in: ui.volatile.activeChannel.dataset.uuid,
                message: this.value,
            }));
            this.value = "";
        }
    });
})();
