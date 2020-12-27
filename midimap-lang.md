# midimap-lang
midimap-lang is the programming language used by MIDIMAP(1) to decide which MIDI events map to which actions.

It consists of a series of COMMENTs and MAPPINGs.

Every line which starts with a hash "#" is a COMMENT. COMMENTs have no effect on the behaviour of the program. Their primary purpose is to describe the non-COMMENT lines, in human-readable language, to other programmers.
```midimap-lang
# This line is a COMMENT.
#The same goes for this line.
 #This line however is not a COMMENT because there is a space preceeding the hash.
```

**Every** line which is not a COMMENT is a MAPPING, there are **no** exceptions to this rule. Although a line is a MAPPING it can nonetheless be an invalid MAPPING.
```midimap-lang
# The empty line following this COMMENT is an invalid MAPPING, but a MAPPING nonetheless.

# The line following this sequence of COMMENTs is an example of a valid MAPPING.
# It is not expected that you understand what the MAPPING does at the moment.
# It is included merely to illustrate that there are valid and invalid MAPPINGs.
# The syntax of a MAPPING will be covered in the next paragraph
part1=20&part2<20->0
```

A valid MAPPING consists of a MATCHER, SEPARATOR and a KEYCODE. The KEYCODE is a integer which corresponds to a keyboard key, the MATCHER is an expression which decides which MIDI events should simulate pressing the KEYCODE and the SEPARATOR is two characters "->" which separate the MATCHER and the KEYCODE. The specific syntax of the MATCHER and KEYCODE components are covered in their own paragraphs below. For readabilities sake, spaces may be inserted **anywhere** in a MAPPING, including the MATCHER, SEPARATOR and KEYCODE components, without affecting the behaviour of the MAPPING.
```midimap-lang
# The line below contains the valid MAPPING introduced in the previous code block.
part1=20&part2<20->0
# As specified a valid MAPPING consists of three components a MATCHER, a SEPARATOR and a KEYCODE.
# We know that the SEPARATOR is the two characters "->" and that it separates the MATCHER and KEYCODE.
# Therefore the MATCHER of this MAPPING is the part preceeding the "->" "part1=20&part2<20" and the KEYCODE is the part following it "0".
# To better illustrate that the MAPPING, SEPARATOR and KEYCODE are distinct components, we will insert a space between each of them
part1=20&part2<20 -> 0
# Since spaces have no affect on the behaviour of the MAPPING, we can technically insert any number of spaces anywhere in the MAPPING without changing the behaviour.
# For example we can insert a space between every character.
p a r t 1 = 2 0 & p a r t 2 < 2 0 - > 0
# The behaviour of the MAPPING above is equal to that of the other MAPPINGs in this codeblock.
# For readability reasons it is best to avoid excessive spacing like in the last example and scant spacing like in the first.
```
