const { Sequelize } = require('sequelize');
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

module.exports = {
  sequelize,
  Product,
  Service,
  About,
  Contact
}; 