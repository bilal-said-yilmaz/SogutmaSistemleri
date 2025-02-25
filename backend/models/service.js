const { DataTypes } = require('sequelize');

module.exports = (sequelize) => {
  return sequelize.define('Service', {
    id: {
      type: DataTypes.UUID,
      defaultValue: DataTypes.UUIDV4,
      primaryKey: true
    },
    title: {
      type: DataTypes.STRING,
      allowNull: false
    },
    description: {
      type: DataTypes.TEXT
    },
    image: {
      type: DataTypes.STRING
    }
  });
}; 