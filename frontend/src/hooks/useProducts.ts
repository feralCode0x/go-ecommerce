import { useState, useEffect } from 'react';
import api from '@utils/axios'; // Tu instancia con interceptor

export const useProducts = () => {
  const [products, setProducts] = useState([]);

  const loadProducts = async () => {
    const response = await api.get('/products');
    setProducts(payload);
  };

  return { products, loadProducts };
};