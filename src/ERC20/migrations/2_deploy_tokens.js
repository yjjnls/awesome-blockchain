const ERC20 = artifacts.require("ERC20");

module.exports = (deployer) => {
    deployer.deploy(ERC20, 10000, 'Simon Bucks', 1, 'SBX');
};
