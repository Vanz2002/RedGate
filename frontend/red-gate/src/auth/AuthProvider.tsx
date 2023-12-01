import React, { useContext, useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useCookies } from 'react-cookie';

const AuthContext = React.createContext<any>('');

export const useAuth = () => {
    return useContext(AuthContext);
}

export const AuthProvider: React.FC<any> = ({ children }) => {
    const [loading, setLoading] = useState(false); // You can set this to false since we assume the user is logged in if the cookie exists
    const [user, setUser] = useState<any>(null);
    const navigate = useNavigate();
    const [cookies] = useCookies(['userID']);

    useEffect(() => {
        const storedUserID = cookies.userID;
        console.log(storedUserID)
        if (storedUserID) {
            setUser(storedUserID);
            console.log("User is already logged in.");
            navigate("/");
        } else {
            console.log("User is not logged in.");
        }

        setLoading(false);
    }, [navigate, cookies]);

    const value = { user }; // You can provide your custom signIn and signOut functions

    return (
        <AuthContext.Provider value={value}>
            {!loading && children}
        </AuthContext.Provider>
    );
}
