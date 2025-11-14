import { useNavigate } from "react-router-dom";
import { LayoutDashboard, Dumbbell, ClipboardList, MapPin, CreditCard, DollarSign, Settings, LogOut, LogIn } from "lucide-react";
import "../styles/Header.css";


const Header = ( ) => {
    const isLoggedIn = localStorage.getItem("isLoggedIn") === "true";
    const isAdmin = localStorage.getItem("isAdmin") === "true";
    const navigate = useNavigate();
    const logout = () => {
        localStorage.removeItem("isLoggedIn");
        localStorage.removeItem("isAdmin");
        localStorage.removeItem("access_token");
        localStorage.removeItem("idUsuario");
        navigate("/");
    }

    return (
        <header>
            <div className="header-container">
                <nav className="header-content">
                    <h1 className="header-title" onClick={() => navigate("/")}>GymPro</h1>
                    <div className="header-links">
                        {isLoggedIn && !isAdmin && (
                            <a href="/dashboard" className="nav-link">
                                <LayoutDashboard size={20} />
                                Dashboard
                            </a>
                        )}
                        <a href="/actividades" className="nav-link">
                            <Dumbbell size={20} />
                            Actividades
                        </a>
                        <a href="/planes" className="nav-link">
                            <ClipboardList size={20} />
                            Planes
                        </a>
                        <a href="/sucursales" className="nav-link">
                            <MapPin size={20} />
                            Sucursales
                        </a>
                        {isLoggedIn && !isAdmin && (
                            <>
                                <a href="/mi-suscripcion" className="nav-link">
                                    <CreditCard size={20} />
                                    Mi Suscripción
                                </a>
                                <a href="/pagos" className="nav-link">
                                    <DollarSign size={20} />
                                    Pagos
                                </a>
                            </>
                        )}
                        {isAdmin && (
                            <a href="/admin" className="nav-link">
                                <Settings size={20} />
                                Panel Admin
                            </a>
                        )}
                        {isLoggedIn ? (
                            <button onClick={logout} className="nav-logout">
                                <LogOut size={20} />
                                Cerrar sesión
                            </button>
                        ) : (
                            <a href="/login" className="nav-link">
                                <LogIn size={20} />
                                Iniciar Sesión
                            </a>
                        )}
                    </div>
                </nav>
            </div>
        </header>
    );
}

export default Header;