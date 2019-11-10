/**
 * @see https://github.com/nektro/mantle
 * @author Meghan Denny <https://nektro.net>
 */
//
(async function() {
    //
    await fetch("/api/about").then(x => x.json()).then(x => {
        el_2.innerText = x.name;
    });
})();
