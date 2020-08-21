// Generated automatically by nearley, version 2.19.6
// http://github.com/Hardmath123/nearley
(function () {
function id(x) { return x[0]; }
var grammar = {
    Lexer: undefined,
    ParserRules: [
    {"name": "main", "symbols": ["statement"], "postprocess": id},
    {"name": "LB", "symbols": [{"literal":"("}]},
    {"name": "RB", "symbols": [{"literal":")"}]},
    {"name": "_", "symbols": []},
    {"name": "_", "symbols": ["_", /[\s]/], "postprocess": function() {}},
    {"name": "statement", "symbols": ["contractFunc"], "postprocess": id},
    {"name": "statement$string$1", "symbols": [{"literal":"a"}, {"literal":"b"}, {"literal":"i"}], "postprocess": function joiner(d) {return d.join('');}},
    {"name": "statement", "symbols": ["statement$string$1"], "postprocess": id},
    {"name": "identifier$ebnf$1", "symbols": []},
    {"name": "identifier$ebnf$1", "symbols": ["identifier$ebnf$1", /[a-z | A-Z | 0-9 | _]/], "postprocess": function arrpush(d) {return d[0].concat([d[1]]);}},
    {"name": "identifier", "symbols": [/[a-z | A-Z]/, "identifier$ebnf$1"], "postprocess": function(d) { return d[0] + d[1].join(""); }},
    {"name": "identifier$ebnf$2", "symbols": [/[ a-z | A-Z | 0-9 | _]/]},
    {"name": "identifier$ebnf$2", "symbols": ["identifier$ebnf$2", /[ a-z | A-Z | 0-9 | _]/], "postprocess": function arrpush(d) {return d[0].concat([d[1]]);}},
    {"name": "identifier", "symbols": [{"literal":"_"}, "identifier$ebnf$2"], "postprocess": function(d) { return d[0] + d[1].join(""); }},
    {"name": "posint", "symbols": [/[0-9]/], "postprocess": id},
    {"name": "posint", "symbols": ["posint", /[0-9]/], "postprocess": function(d) { return d[0] + d[1]; }},
    {"name": "int", "symbols": [{"literal":"-"}, "posint"], "postprocess": function(d) { return d[0] + d[1]; }},
    {"name": "int", "symbols": ["posint"], "postprocess": id},
    {"name": "float", "symbols": ["int"], "postprocess": id},
    {"name": "float", "symbols": ["int", {"literal":"."}, "posint"], "postprocess": function(d) {return d[0] + d[1] + d[2]; }},
    {"name": "string", "symbols": [{"literal":"\""}, "_string", {"literal":"\""}], "postprocess": function(d) {return d[1]; }},
    {"name": "string", "symbols": [{"literal":"'"}, "_sString", {"literal":"'"}], "postprocess": function(d) {return d[1]; }},
    {"name": "_sString", "symbols": [], "postprocess": function() {return ""; }},
    {"name": "_sString", "symbols": ["_sString", "_sStringchar"], "postprocess": function(d) {return d[0] + d[1];}},
    {"name": "_sStringchar", "symbols": [/[^\\']/], "postprocess": id},
    {"name": "_sStringchar", "symbols": [{"literal":"\\"}, /[^]/], "postprocess": function(d) {return JSON.parse("\"" + d[0] + d[1] + "\""); }},
    {"name": "_string", "symbols": [], "postprocess": function() {return ""; }},
    {"name": "_string", "symbols": ["_string", "_stringchar"], "postprocess": function(d) {return d[0] + d[1];}},
    {"name": "_stringchar", "symbols": [/[^\\"]/], "postprocess": id},
    {"name": "_stringchar", "symbols": [{"literal":"\\"}, /[^]/], "postprocess": function(d) {return JSON.parse("\"" + d[0] + d[1] + "\""); }},
    {"name": "_", "symbols": []},
    {"name": "_", "symbols": ["_", /[\s]/], "postprocess": function() {}},
    {"name": "__", "symbols": [/[\s]/]},
    {"name": "__", "symbols": ["__", /[\s]/], "postprocess": function() {}},
    {"name": "fileSpecifier$string$1", "symbols": [{"literal":"f"}, {"literal":":"}], "postprocess": function joiner(d) {return d.join('');}},
    {"name": "fileSpecifier", "symbols": ["fileSpecifier$string$1", "string"], "postprocess": function(d) {return {_type: 'file', value: d[1]}; }},
    {"name": "arg$string$1", "symbols": [{"literal":"t"}, {"literal":"r"}, {"literal":"u"}, {"literal":"e"}], "postprocess": function joiner(d) {return d.join('');}},
    {"name": "arg", "symbols": ["arg$string$1"], "postprocess": function(d) { return true; }},
    {"name": "arg$string$2", "symbols": [{"literal":"f"}, {"literal":"a"}, {"literal":"l"}, {"literal":"s"}, {"literal":"e"}], "postprocess": function joiner(d) {return d.join('');}},
    {"name": "arg", "symbols": ["arg$string$2"], "postprocess": function(d) { return false; }},
    {"name": "arg", "symbols": ["fileSpecifier"], "postprocess": id},
    {"name": "arg", "symbols": ["int"], "postprocess": function(d) { return Number(d[0]); }},
    {"name": "arg", "symbols": ["float"], "postprocess": function(d) { return Number(d[0]); }},
    {"name": "arg", "symbols": ["string"], "postprocess": id},
    {"name": "argList", "symbols": ["_", "arg", "_"], "postprocess": function(d) { return [d[1]]; }},
    {"name": "argList", "symbols": ["argList", "_", {"literal":","}, "_", "arg"], "postprocess": function(d) {  return d[0].concat(d[4]); }},
    {"name": "args", "symbols": ["LB", "RB"], "postprocess": function(d) { return []; }},
    {"name": "args", "symbols": ["LB", "argList", "RB"], "postprocess": function(d) { return d[1]; }},
    {"name": "contractFunc", "symbols": ["identifier", "args"], "postprocess": function(d) { return { name: d[0], args: d[1] }; }}
]
  , ParserStart: "main"
}
if (typeof module !== 'undefined'&& typeof module.exports !== 'undefined') {
   module.exports = grammar;
} else {
   window.grammar = grammar;
}
})();
