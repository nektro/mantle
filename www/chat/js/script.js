"use strict";
//
import "./x/index.js";
//
import { create_element, dcTN, setDataBinding } from "./util.js";
import * as ui from "./ui.js";
import { el_2, el_3, el_1, output, el_uonline, el_input, context, audio_buffer_size } from "./ui.util.js";
import * as api from "./api/index.js";
import * as ws from "./ws.js";

//
$(document).on("click", (e) => {
    const p = e.target.path();
    if (p.filter((v) => v.matches("dialog[open]")).length > 0) return;
    if (p.filter((v) => v.matches(".msg .usr")).length > 0) return;
    const s = document.querySelectorAll("dialog[open]");
    s.forEach((v) => v.removeAttribute("open"));
});
window.addEventListener("blur", () => {
    ui.volatile.windowActive = false;
});
window.addEventListener("focus", () => {
    ui.volatile.windowActive = true;
});
document.getElementById("shrink_uonline").addEventListener("click", () => {
    output.classList.toggle("extended-right");
    el_uonline.toggleAttribute("hidden");
    el_input.classList.toggle("extended-right");
    output.children[0]._scroll_to_bottom();
});

//
(async function() {
    //
    moment.defaultFormat = "ddd MMM DD Y HH:mm:ss zZZ";

    //
    await api.M.meta.about().then((x) => {
        setDataBinding("server_name", x.name);
        setDataBinding("server_version", x.version);
        //
        const sx = document.querySelector("x-settings[data-s-for=server] [data-s-section=overview]");
        const sxc = sx.querySelectorAll("[endpoint][name]");
        for (const item of sxc) {
            const n = item.getAttribute("name");
            sx.querySelector(`[name="${n}"]`).setAttribute("value", x[n]);
        }
        //
        el_2.children[1].addEventListener("click", () => {
            document.querySelector("x-settings[data-s-for=server]")._open();
        });
    });

    //
    await api.M.users.me().then((x) => {
        ui.volatile.me = x.user;
        ui.volatile.me.perms = x.perms;
        const n = ui.volatile.me.nickname || ui.volatile.me.name;
        el_3.children[0].textContent = `@${n}`;
        document.querySelectorAll("[data-requires]").forEach((el) => {
            el.setAttribute("hidden", "");
        });
        const p = x.perms;
        for (const key in p) {
            if (p[key]) {
                document.querySelectorAll(`[data-requires^="${key}"]`).forEach((el) => {
                    el.removeAttribute("hidden");
                });
            }
        }
        //
        el_3.children[1].addEventListener("click", () => {
            document.querySelector("x-settings[data-s-for=user]")._open();
        });
        //
        document.querySelectorAll("[data-s-for='user'] [data-s-section='my_account'] [fill]").forEach((el) => {
            el.setAttribute("fill", x.user.uuid);
        });
        document.querySelector("[data-s-for='user'] [data-s-section='my_account'] [name='nickname']").setAttribute("value", x.user.nickname);
        //
        if (Notification.permission !== "granted") {
            ui.toggleHandlers.get("notifications_messages")("0");
        }
    }).catch(() => {
        location.assign("../");
    });

    //
    await api.M.roles.me().then((x) => {
        const rls = x.sort((a,b) => a.position > b.position);
        //
        for (const item of rls) {
            ui.M.role.add(item);
        }
    });

    //
    await api.M.channels.me().then(async (x) => {
        for (const item of x) {
            await ui.M.channel.add(item);
        }
        await output.setActiveChannel(x[0].uuid);

        el_1.querySelector("button").addEventListener("click", async () => {
            const {value: name} = await Swal({
                title: "Enter the new channel's name",
                input: "text",
                showCancelButton: true,
                inputValidator: (value) => !value && "You need to write something!",
            });
            if (name !== undefined) {
                return api.M.channels.create(name);
            }
        });
        el_1.querySelector("ol").addEventListener("click", (ev) => {
            const fl = ev.composedPath().filter((v) => v instanceof Element && v.matches("[data-uuid]"));
            if (fl.length === 0) return;
            output.setActiveChannel(fl[0].dataset.uuid);
        });
    });

    await api.M.invites.me().then((x) => {
        for (const item of x) {
            ui.M.invite.add(item);
        }
    });

    await api.M.users.online().then((x) => {
        for (const item of x) {
            ui.M.user.connect(item);
        }
    });

    //
    /** @type {HTMLInputElement} */
    const input = document.getElementById("input").children[0];
    const socket = new WebSocket(`ws${location.protocol.substring(4)}//${location.host}/ws`);

    socket.addEventListener("open", () => {
        el_2.children[0].classList.remove("loading");
        el_2.children[0].classList.add("online");
    });
    socket.addEventListener("close", () => {
        el_2.children[0].classList.remove("online");
    });
    socket.addEventListener("message", async (e) => {
        const d = JSON.parse(e.data);
        let o = ws.M;
        for (const item of d.type.split("-")) {
            if (!(item in o)) {
                console.error("event handler not found:", d);
                return;
            }
            o = o[item];
        }
        if (typeof o !== "function") {
            console.error("handler is not a function:", `"${d.type}"`, `"${typeof o}"`);
            return;
        }
        await o(d);
    });
    setInterval(() => {
        if (el_2.children[0].classList.contains("online")) {
            socket.send(JSON.stringify({
                type: "ping",
            }));
        }
    }, 30*1000);

    //
    input.addEventListener("keydown", (e) => {
        if (e.key === "Enter") {
            const msg_con = e.target.value;
            if (msg_con.length === 0) return;
            socket.send(JSON.stringify({
                type: "message",
                in: output.active_channel_uid,
                message: msg_con.trim(),
            }));
            e.target.value = "";
        }
    });

    let voice_connected = false;
    /** @type {MediaStreamAudioSourceNode} */
    let voice_source = null;

    //
    document.getElementById("voice_chat").addEventListener("click", () => {
        const usr_list = document.getElementById("voice_chat").children[1];
        const usr_name = ui.volatile.me.getName() + "#" + ui.volatile.me.id;
        const usr_uuid = ui.volatile.me.uuid;

        if (!voice_connected) {
            navigator.mediaDevices.getUserMedia({ audio: true }).then((stream) => {
                voice_connected = true;
                console.debug("connected");
                voice_source = context.createMediaStreamSource(stream);
                const source = context.createMediaStreamSource(stream);
                const processor = context.createScriptProcessor(audio_buffer_size, 1, 1);
                source.connect(processor);
                processor.connect(context.destination);

                processor.addEventListener("audioprocess", (e) => {
                    if (!voice_connected) { return; }
                    const data = [...e.inputBuffer.getChannelData(0)];
                    socket.send(JSON.stringify({
                        type: "voice-data",
                        from: ui.volatile.me.uuid,
                        data: data,
                    }));
                });

                usr_list.appendChild(create_element("li", [["data-uuid",usr_uuid]], [dcTN(usr_name)]));
            });
        }
        else {
            voice_connected = false;
            console.debug("disconnected");
            voice_source.disconnect();
            Array.from(usr_list.children).forEach((v) => {
                if (v.dataset.uuid === usr_uuid) {
                    v.remove();
                }
            });
        }
    });

    //
    console.log("finished init");
})();
