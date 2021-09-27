pragma solidity ^0.8.7;

interface ERC165 {
    /// @notice Query if a contract implements an interface
    /// @param interfaceID The interface identifier, as specified in ERC-165
    /// @dev Interface identification is specified in ERC-165. This function
    ///  uses less than 30,000 gas.
    /// @return `true` if the contract implements `interfaceID` and
    ///  `interfaceID` is not 0xffffffff, `false` otherwise
    function supportsInterface(bytes4 interfaceID) external view returns (bool);
}

/// @title ERC-173 Contract Ownership Standard
///  Note: the ERC-165 identifier for this interface is 0x7f5828d0
interface ERC173 is ERC165 {
    /// @dev This emits when ownership of a contract changes.
    event OwnershipTransferred(
        address indexed previousOwner,
        address indexed newOwner
    );

    /// @notice Get the address of the owner
    /// @return The address of the owner.
    function owner() external view returns (address);

    /// @notice Set the address of the new owner of the contract
    /// @dev Set _newOwner to address(0) to renounce any ownership.
    /// @param _newOwner The address of the new owner of the contract
    function transferOwnership(address _newOwner) external;
}

contract Register is ERC173 {
    uint256 public constant PRICE = 1 ether / 10; // Registration fee = 0.1 ether
    address private _owner;

    mapping(address => bool) _registered;

    constructor() {
        _owner = msg.sender;
    }

    modifier onlyOwner() {
        require(
            msg.sender == _owner,
            "Must be the owner to perform this operation"
        );
        _;
    }

    /// @notice Get the address of the owner
    /// @return The address of the owner.
    function owner() external view override returns (address) {
        return _owner;
    }

    /// @notice Set the address of the new owner of the contract
    /// @dev Set _newOwner to address(0) to renounce any ownership.
    /// @param _newOwner The address of the new owner of the contract
    function transferOwnership(address _newOwner) external override onlyOwner {
        address old = _owner;
        _owner = _newOwner;
        emit OwnershipTransferred(old, _newOwner);
    }

    /// @notice Query if a contract implements an interface
    /// @param interfaceID The interface identifier, as specified in ERC-165
    /// @dev Interface identification is specified in ERC-165. This function
    ///  uses less than 30,000 gas.
    /// @return `true` if the contract implements `interfaceID` and
    ///  `interfaceID` is not 0xffffffff, `false` otherwise
    function supportsInterface(bytes4 interfaceID)
        external
        pure
        override
        returns (bool)
    {
        return
            interfaceID == this.supportsInterface.selector ||
            interfaceID ==
            this.owner.selector ^ this.transferOwnership.selector;
    }

    function isRegistered(address _address) external view returns (bool) {
        return _registered[_address];
    }

    function register() external payable {
        require(!this.isRegistered(msg.sender), "You are already registered!");
        require(msg.value == PRICE, "Registration price is 0.1 ETH!");

        
        _registered[msg.sender] = true;
    }

    function withdraw() external onlyOwner {
        payable(_owner).transfer(address(this).balance);
    }
}
