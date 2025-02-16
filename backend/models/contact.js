const { DataTypes } = require('sequelize');

module.exports = (sequelize) => {
  return sequelize.define('Contact', {
    id: {
      type: DataTypes.UUID,
      defaultValue: DataTypes.UUIDV4,
      primaryKey: true
    },
    title: {
      type: DataTypes.STRING,
      allowNull: false
    },
    phone: {
      type: DataTypes.STRING
    },
    email: {
      type: DataTypes.STRING
    },
    address: {
      type: DataTypes.TEXT
    }
  });
}; 