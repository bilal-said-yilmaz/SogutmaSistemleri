const { Sequelize, DataTypes } = require('sequelize');
const path = require('path');
require('dotenv').config();

const sequelize = new Sequelize(process.env.DATABASE_URL, {
  dialect: 'postgres',
  logging: false
});

const Product = require('./product')(sequelize);
const Service = require('./service')(sequelize);
const About = require('./about')(sequelize);
const Contact = require('./contact')(sequelize);

const Hero = sequelize.define('Hero', {
  subheading: {
    type: DataTypes.STRING,
    allowNull: false
  },
  heading: {
    type: DataTypes.STRING,
    allowNull: false
  },
  buttonText: {
    type: DataTypes.STRING,
    allowNull: false
  },
  backgroundImage: {
    type: DataTypes.TEXT,
    allowNull: true
  }
})

const Footer = sequelize.define('Footer', {
  copyright: {
    type: DataTypes.STRING,
    allowNull: false
  },
  socialLinks: {
    type: DataTypes.JSON,
    allowNull: false
  },
  links: {
    type: DataTypes.JSON,
    allowNull: false
  }
})

module.exports = {
  sequelize,
  Product,
  Service,
  About,
  Contact,
  Hero,
  Footer
}; 