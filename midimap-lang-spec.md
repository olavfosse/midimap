# midimap-lang specification
## Writing style
This document consists of a series of definitions of named syntactic elements. To prevent confusion between a named syntactic element and similar sounding words, we always refer to syntactic elements with all uppercase. Note that not all uppercase words are referring to syntactic elements, an example being MIDI which refers to the MIDI protocol.
## 1 COMMENTS and MAPPINGS
A midimap-lang program consists of a series of COMMENTS and MAPPINGS. Each line is either a COMMENT or a MAPPING and each COMMENT or MAPPING is contained entirely within this line.
A COMMENT is a line which has no impact on the behaviour of the program. Any line where the first character is a hash "#" is a COMMENT.
### 1.2 MAPPINGS
A MAPPING is a line which instructs the program to map incoming MIDI events to a simulated keypress with *some* keycode if *some* criteria is met. Every line which is not a COMMENT is a MAPPING.

A MAPPING is of the form "MATCHER -> KEYCODE". The KEYCODE specifies the keycode of the keypress that will be simulated if the criteria met. The MATCHER specifies the criteria.
