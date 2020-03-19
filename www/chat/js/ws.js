"use strict";
//
import * as ui from "./ui.js";
import * as api from "./api/index.js";
import { Channel } from "./ui.channel.js";
import { output, messageCache } from "./ui.util.js";
import { setDataBinding } from "./util.js";

//
export const M = {
    pong: () => {
        // do nothing, keep connection alive
    },
    user: {
        connect: (d) => {
            ui.setMemberOnline(d.user);
        },
        disconnect: (d) => {
            ui.setMemberOffline(d.user);
        },
        update: (d) => {
            new api.User(d.user);
            if (["add_role","remove_role"].includes(d.key)) {
                document.querySelectorAll(`dialog.popup.user ol [data-role="${d.value}"]`).forEach((v) => {
                    v.classList.toggle("active");
                });
            }
        },
    },
    channel: {
        new: (d) => {
            const h = new api.Channel(d.channel);
            ui.addChannel(h);
        },
        update: (d) => {
            const c = new api.Channel(d.channel);
            if (["name"].includes(d.key)) {
                new Channel(c.uuid).p_name = d.value;
            }
            if (output.dataset.active === c.uuid) {
                setDataBinding("channel_name", c.name);
                setDataBinding("channel_description", c.description);
            }
        },
    },
    role: {
        new: (d) => {
            const r = new api.Role(d.role);
            ui.addRole(r);
        },
        update: (d) => {
            new api.Role(d.role);
            if (["color"].includes(d.key)) {
                const x = document.getElementById("link-role-color");
                const y = x.href.split("=");
                x.href = y[0]+"="+(parseInt(y[1],10)+1).toString();
            }
        },
    },
    message: {
        new: async (d) => {
            const u = await api.M.users.get(d.message.author);
            await ui.addMessage(d.in, u, d.message, true);
        },
        delete: (d) => {
            const a1 = messageCache.get(d.channel);
            const a2 = a1.filter((v) => d.affected.includes(v.uuid));
            a2.forEach((v) => a1.splice(a1.indexOf(v), 1));
            for (const item of d.affected) {
                const de = output.querySelector(`.msg[data-msg-uid="${item}"]`);
                if (de !== null) de.remove();
            }
        },
    },
    invite: {
        new: (d) => {
            ui.addInvite(new api.Invite(d.invite));
        },
    },
};
