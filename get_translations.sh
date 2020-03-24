#!/usr/bin/env bash

project="390627"

_crowdin_get() {
    curl -s -X GET -H "Accept: application/json" -H "Authorization: Bearer $CROWDIN_ACCESS_TOKEN" "https://api.crowdin.com/api/v2/projects/$project$1"
}
_crowdin_jq_c() {
    _crowdin_get "$1" | jq -c "$2"
}
_crowdin_jq_r() {
    _crowdin_get "$1" | jq -cr "$2"
}
_jq_r() {
    echo "$1" | jq -cr "$2"
}

#
mkdir -p 'www_bin/translations'
fileM="www_bin/translations/_languages.json"
printf "[" > "$fileM"
#
_crowdin_jq_r '' '.data.targetLanguageIds[]' |
while IFS= read -r langID; do
    echo "$langID:"
    printf "\"$langID\"," >> "$fileM"
    fileL="www_bin/translations/$langID.json"
    echo "{" > "$fileL"
    #
    _crowdin_jq_c '/strings?limit=500' '.data[] | .data | {id,context,text}' |
    while IFS= read -r strJSN; do
        strID=$(_jq_r "$strJSN" '.id')
        strCN=$(_jq_r "$strJSN" '.context')
        strTX=$(_jq_r "$strJSN" '.text')
        #
        translation=$(_crowdin_jq_r "/translations?limit=500&languageId=$langID&stringId=$strID" '.data | sort_by(.data.rating)[-1].data.text | select (.!=null)')
        result=${translation:-$strTX}
        line2="\"$strCN\":\"$result\","
        echo "$langID: $strCN: $result"
        echo "$line2" >> "$fileL"
    done
    echo '"____":""' >> "$fileL"
    echo "}" >> "$fileL"
done
echo '""]' >> "$fileM"
