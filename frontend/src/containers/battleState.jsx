import {connect} from 'react-redux'


const getBattleState = () => {

}

const mapStateToProps = state => ({
	battleState: getBattleState()
})

const mapDispatchToProps = dispatch => ({
	placeShip: (id, location) => dispatch(placeShip(id, location))
})

export default connect(mapStateToProps, mapDispatchToProps)(BattleState)
