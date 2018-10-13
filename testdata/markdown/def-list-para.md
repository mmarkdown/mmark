Option Type
: 8-bit identifier of the type of option.  The option identifier

    The highest-order two bits are set to 00 so any node not
    implementing Scenic Routing will skip over this option and
    continue processing the header.  The third-highest-order bit

    indicates that the SRO does not change en route to the packet's
    final destination.

Option Length
: 8-bit unsigned integer.  The length of the option in octets
