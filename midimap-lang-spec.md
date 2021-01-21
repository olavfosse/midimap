# midimap-lang specification
## Writing style
This document consists of a series of definitions of named syntactic elements. To prevent confusion between a named syntactic element and similar sounding words, we always refer to syntactic elements with all uppercase. Note that not all uppercase words are referring to syntactic elements, an example being MIDI which refers to the MIDI protocol.
## 1 COMMENTS and MAPPINGS
A midimap-lang program consists of a series of COMMENTS and MAPPINGS. Each line is either a COMMENT or a MAPPING and each COMMENT or MAPPING is contained entirely within this line.
### 1.1 COMMENTS
A COMMENT is a line which has no impact on the behaviour of the program. Every line where the first character is a hash "#" is a COMMENT.
### 1.2 MAPPINGS
A MAPPING is a line which instructs the program to map incoming MIDI events to a simulated keypress with *some* keycode if *some* criteria is met. Every line which is not a COMMENT is a MAPPING.

A MAPPING is of the form "MATCHER -> KEYCODE".
#### 1.2.1 MATCHERS
A MATCHER specifies the criteria required by its parent MAPPING to simulate a keypress.

A MATCHER is of the form "LEFTOPERAND COMPARISONOPERATOR RIGHTOPERAND" or "MATCHER LOGICALOPERATOR MATCHER".

LEFTOPERAND is "data1" or "data2" representing the first or second data byte, COMPARISONOPERATOR is "<", "<=", "==", "!=", ">=" or ">" representing a comparison operator and RIGHTOPERAND is a series of decimal digits representing an integer.

LOGICALOPERATOR is "||" or "&&" representing a logical operator.
#### 1.2.2 KEYCODES
A KEYCODE specifies the keycode which its parent MAPPING will simulate a keypress with if the criteria specified by its MATCHER is fulfilled.

A KEYCODE is an integer represented by a series of decimal digits.
