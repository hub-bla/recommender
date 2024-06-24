import './Movie.css'
interface MovieProps {
    idx: Number
    title: string
}

export const Movie: React.FC<MovieProps> = ({idx, title}) => {
    return (
        <div className='movie-component' onClick={() => console.log(idx)}>
            <h3>{title}</h3>
        </div>
    )
}