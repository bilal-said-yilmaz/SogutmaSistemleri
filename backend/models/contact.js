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
    },
    latitude: {
      type: DataTypes.FLOAT,
      defaultValue: 37.4538
    },
    longitude: {
      type: DataTypes.FLOAT,
      defaultValue: 35.8158
    },
    weekdayHours: {
      type: DataTypes.STRING,
      defaultValue: '09:00 - 18:00'
    },
    saturdayHours: {
      type: DataTypes.STRING,
      defaultValue: '09:00 - 14:00'
    },
    sundayHours: {
      type: DataTypes.STRING,
      defaultValue: 'Kapalı'
    }
  });
}; 